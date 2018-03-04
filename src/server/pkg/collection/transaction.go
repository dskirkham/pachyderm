package collection

// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// We copy this code from etcd because the etcd implementation of STM does
// not have the DelAll method, which we need.

import (
	"fmt"

	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
)

// STM is an interface for software transactional memory.
type STM interface {
	// Get returns the value for a key and inserts the key in the txn's read set.
	// If Get fails, it aborts the transaction with an error, never returning.
	Get(key string) (string, error)
	// Put adds a value for a key to the write set.
	Put(key, val string, opts ...v3.OpOption)
	// Rev returns the revision of a key in the read set.
	Rev(key string) int64
	// Del deletes a key.
	Del(key string)
	// DelAll deletes all keys with the given prefix
	// Note that the current implementation of DelAll is incomplete.
	// To use DelAll safely, do not issue any Get/Put operations after
	// DelAll is called.
	DelAll(key string)
	Context() context.Context

	// commit attempts to apply the txn's changes to the server.
	commit() *v3.TxnResponse
	reset()
}

// stmError safely passes STM errors through panic to the STM error channel.
type stmError struct{ err error }

// NewSTM intiates a new STM operation. It uses a serializable model.
func NewSTM(ctx context.Context, c *v3.Client, apply func(STM) error) (*v3.TxnResponse, error) {
	return newSTMSerializable(ctx, c, apply)
}

// newSTMRepeatable initiates new repeatable read transaction; reads within
// the same transaction attempt always return the same data.
func newSTMRepeatable(ctx context.Context, c *v3.Client, apply func(STM) error) (*v3.TxnResponse, error) {
	s := &stm{client: c, ctx: ctx, getOpts: []v3.OpOption{v3.WithSerializable()}}
	return runSTM(s, apply)
}

// newSTMSerializable initiates a new serialized transaction; reads within the
// same transactiona attempt return data from the revision of the first read.
func newSTMSerializable(ctx context.Context, c *v3.Client, apply func(STM) error) (*v3.TxnResponse, error) {
	s := &stmSerializable{
		stm:      stm{client: c, ctx: ctx},
		prefetch: make(map[string]*v3.GetResponse),
	}
	return runSTM(s, apply)
}

// newSTMReadCommitted initiates a new read committed transaction.
func newSTMReadCommitted(ctx context.Context, c *v3.Client, apply func(STM) error) (*v3.TxnResponse, error) {
	s := &stmReadCommitted{stm{client: c, ctx: ctx, getOpts: []v3.OpOption{v3.WithSerializable()}}}
	return runSTM(s, apply)
}

type stmResponse struct {
	resp *v3.TxnResponse
	err  error
}

func runSTM(s STM, apply func(STM) error) (*v3.TxnResponse, error) {
	outc := make(chan stmResponse, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e, ok := r.(stmError)
				if !ok {
					// client apply panicked
					panic(r)
				}
				outc <- stmResponse{nil, e.err}
			}
		}()
		var out stmResponse
		for {
			s.reset()
			if out.err = apply(s); out.err != nil {
				break
			}
			if out.resp = s.commit(); out.resp != nil {
				break
			}
		}
		outc <- out
	}()
	r := <-outc
	return r.resp, r.err
}

// stm implements repeatable-read software transactional memory over etcd
type stm struct {
	client *v3.Client
	ctx    context.Context
	// rset holds read key values and revisions
	rset map[string]*v3.GetResponse
	// wset holds overwritten keys and their values
	wset map[string]stmPut
	// getOpts are the opts used for gets
	getOpts []v3.OpOption
}

type stmPut struct {
	val string
	op  v3.Op
}

func (s *stm) Context() context.Context {
	return s.ctx
}

func (s *stm) Get(key string) (string, error) {
	if wv, ok := s.wset[key]; ok {
		return wv.val, nil
	}
	return respToValue(key, s.fetch(key))
}

func (s *stm) Put(key, val string, opts ...v3.OpOption) {
	s.wset[key] = stmPut{val, v3.OpPut(key, val, opts...)}
}

func (s *stm) Del(key string) { s.wset[key] = stmPut{"", v3.OpDelete(key)} }

func (s *stm) DelAll(key string) { s.wset[key] = stmPut{"", v3.OpDelete(key, v3.WithPrefix())} }

func (s *stm) Rev(key string) int64 {
	if resp := s.fetch(key); resp != nil && len(resp.Kvs) != 0 {
		return resp.Kvs[0].ModRevision
	}
	return 0
}

func (s *stm) commit() *v3.TxnResponse {
	cmps := s.cmps()
	puts := s.puts()
	txnresp, err := s.client.Txn(s.ctx).If(cmps...).Then(puts...).Commit()
	if err == rpctypes.ErrTooManyOps {
		panic(stmError{
			fmt.Errorf(
				"%v (%d comparisons, %d puts: hint: set --max-txn-ops on the "+
					"ETCD cluster to at least the largest of those values)",
				err, len(cmps), len(puts)),
		})
	} else if err != nil {
		panic(stmError{err})
	}
	if txnresp.Succeeded {
		return txnresp
	}
	return nil
}

// cmps guards the txn from updates to read set
func (s *stm) cmps() []v3.Cmp {
	cmps := make([]v3.Cmp, 0, len(s.rset))
	for k, rk := range s.rset {
		cmps = append(cmps, isKeyCurrent(k, rk))
	}
	return cmps
}

func (s *stm) fetch(key string) *v3.GetResponse {
	if resp, ok := s.rset[key]; ok {
		return resp
	}
	resp, err := s.client.Get(s.ctx, key, s.getOpts...)
	if err != nil {
		panic(stmError{err})
	}
	s.rset[key] = resp
	return resp
}

// puts is the list of ops for all pending writes
func (s *stm) puts() []v3.Op {
	puts := make([]v3.Op, 0, len(s.wset))
	for _, v := range s.wset {
		puts = append(puts, v.op)
	}
	return puts
}

func (s *stm) reset() {
	s.rset = make(map[string]*v3.GetResponse)
	s.wset = make(map[string]stmPut)
}

type stmSerializable struct {
	stm
	prefetch map[string]*v3.GetResponse
}

func (s *stmSerializable) Get(key string) (string, error) {
	if wv, ok := s.wset[key]; ok {
		return wv.val, nil
	}
	firstRead := len(s.rset) == 0
	if resp, ok := s.prefetch[key]; ok {
		delete(s.prefetch, key)
		s.rset[key] = resp
	}
	resp := s.stm.fetch(key)
	if firstRead {
		// txn's base revision is defined by the first read
		s.getOpts = []v3.OpOption{
			v3.WithRev(resp.Header.Revision),
			v3.WithSerializable(),
		}
	}
	return respToValue(key, resp)
}

func (s *stmSerializable) Rev(key string) int64 {
	s.Get(key)
	return s.stm.Rev(key)
}

func (s *stmSerializable) gets() ([]string, []v3.Op) {
	keys := make([]string, 0, len(s.rset))
	ops := make([]v3.Op, 0, len(s.rset))
	for k := range s.rset {
		keys = append(keys, k)
		ops = append(ops, v3.OpGet(k))
	}
	return keys, ops
}

func (s *stmSerializable) commit() *v3.TxnResponse {
	keys, getops := s.gets()
	cmps := s.cmps()
	puts := s.puts()
	txn := s.client.Txn(s.ctx).If(cmps...).Then(puts...)
	// use Else to prefetch keys in case of conflict to save a round trip
	txnresp, err := txn.Else(getops...).Commit()
	if err == rpctypes.ErrTooManyOps {
		panic(stmError{
			fmt.Errorf(
				"%v (%d comparisons, %d puts: hint: set --max-txn-ops on the "+
					"ETCD cluster to at least the largest of those values)",
				err, len(cmps), len(puts)),
		})
	} else if err != nil {
		panic(stmError{err})
	}
	if txnresp.Succeeded {
		return txnresp
	}
	// load prefetch with Else data
	for i := range keys {
		resp := txnresp.Responses[i].GetResponseRange()
		s.rset[keys[i]] = (*v3.GetResponse)(resp)
	}
	s.prefetch = s.rset
	s.getOpts = nil
	return nil
}

type stmReadCommitted struct{ stm }

// commit always goes through when read committed
func (s *stmReadCommitted) commit() *v3.TxnResponse {
	s.rset = nil
	return s.stm.commit()
}

func isKeyCurrent(k string, r *v3.GetResponse) v3.Cmp {
	if len(r.Kvs) != 0 {
		return v3.Compare(v3.ModRevision(k), "=", r.Kvs[0].ModRevision)
	}
	return v3.Compare(v3.ModRevision(k), "=", 0)
}

func respToValue(key string, resp *v3.GetResponse) (string, error) {
	if len(resp.Kvs) == 0 {
		return "", ErrNotFound{Key: key}
	}
	return string(resp.Kvs[0].Value), nil
}
