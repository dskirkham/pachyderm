package glob

import "testing"

// tests derived from https://github.com/micromatch/extglob/test
func TestExtGlob(t *testing.T) {
	for _, test := range []test{
		glob(true, "*.(a|b)", "a.a.a"),
		glob(true, "*.(a|b)", "a.a"),
		glob(true, "*.(a|b)", "a.b"),
		glob(false, "*.(a|b)", "a.bb"),
		glob(false, "*.(a|b)", "a.ccc"),
		glob(true, "*.(a|b)", "c.a"),
		glob(false, "*.(a|b)", "d.a.d"),
		glob(true, "*.(a|b)*", "a.a.a"),
		glob(true, "*.(a|b)*", "a.a"),
		glob(true, "*.(a|b)*", "a.b"),
		glob(true, "*.(a|b)*", "a.bb"),
		glob(false, "*.(a|b)*", "a.ccc"),
		glob(true, "*.(a|b)*", "c.a"),
		glob(true, "*.(a|b)*", "d.a.d"),

		glob(false, "??", "a"),
		glob(false, "??", "aab"),
		glob(false, "?", "aa"),
		glob(false, "?", "aab"),
		glob(true, "?", "a"),
		glob(true, "?(a*|b)", "ax"),
		glob(true, "?(a|b)", "a"),
		glob(true, "?(ab)/..?(/)", "ab/../"),
		glob(true, "?(ab|??)/..?(/)", "ab/../"),
		glob(true, "?@(a|b)*@(c)d", "abcd"),
		glob(true, "?b/..?(/)", "ab/../"),
		glob(true, "(*(a|b\\[)|f*)", "foo"),
		glob(true, "(a)", "a"),

		glob(false, "(a+|b)*", "zz"),
		glob(false, "(a+|b)+", "abcdef"),
		glob(false, "(a+|b)+", "abcfef"),
		glob(false, "(a+|b)+", "abcfefg"),
		glob(false, "(a+|b)+", "abd"),
		glob(false, "(a+|b)+", "abef"),
		glob(false, "(a+|b)+", "acd"),
		glob(true, "(a|d).(a|b)*", "a.a"),
		glob(true, "(a|d).(a|b)*", "a.b"),
		glob(true, "(a|d).(a|b)*", "a.bb"),
		glob(false, "(b)", "a"),

		glob(true, "[!/][!/]/../", "ab/../"),
		glob(true, "[^/][^/]/../", "ab/../"),
		glob(false, "[a*\\(]*z", "abcx"),
		glob(false, "[a*\\(]*z", "bbc"),
		glob(true, "[a*\\(]*z", "abcz"),
		glob(true, "@(??)/..?(/)", "ab/../"),
		glob(true, "@(??|a*)/..?(/)", "ab/../"),
		glob(true, "@(?b|?b)/..?(/)", "ab/../"),
		glob(true, "@(a?|?b)/..?(/)", "ab/../"),
		glob(false, "@(a)b", "aa"),
		glob(true, "@(a*)/..?(/)", "ab/../"),
		glob(true, "@(ab|?b)/..?(/)", "ab/../"),
		glob(true, "@(ab|+([!/]))/..?(/)", "ab/../"),
		glob(true, "@(ab|+([^/]))/..?(/)", "ab/../"),
		glob(true, "@(ab|a*(b))*(c)d", "abbcd"),
		glob(true, "@(ab|a*(b))*(c)d", "acd"),
		glob(true, "@(ab|a*@(b))*(c)d", "abcd"),
		glob(true, "@(b+(c)d|e*(f)g?|?(h)i@(j|k))", "effgz"),
		glob(true, "@(b+(c)d|e*(f)g?|?(h)i@(j|k))", "efgz"),
		glob(true, "@(b+(c)d|e*(f)g?|?(h)i@(j|k))", "egz"),
		glob(false, "@(b+(c)d|e+(f)g?|?(h)i@(j|k))", "egz"),
		glob(false, "@(c)b", "aab"),
		glob(true, "@(foo|f|fo)*(f|of+(o))", "foofoofo"),
		glob(true, "@(x)", "x"),

		glob(false, "*;[1-9]*([0-9])", "\"MS.FILE;1\""),
		glob(false, "*;[1-9]*([0-9])", "\"MS.FILE;13\""),
		glob(true, "*;[1-9]*([0-9])*", "\"MS.FILE;13\""),
		glob(true, "*;[1-9]**([0-9])*", "\"MS.FILE;13\""),
		glob(true, "*?(a)bc", "123abc"),

		glob(false, "*.+(b|d)", "a."),
		glob(false, "*.+(b|d)", "a.a.a"),
		glob(false, "*.+(b|d)", "a.a"),
		glob(false, "*.+(b|d)", "a.ccc"),
		glob(false, "*.+(b|d)", "c.a"),
		glob(true, "*.+(b|d)", "a.b"),
		glob(true, "*.+(b|d)", "a.bb"),
		glob(true, "*.+(b|d)", "d.a.d"),
		glob(false, "*.c?(c)", "file.C"),
		glob(false, "*.c?(c)", "file.ccc"),
		glob(true, "*.c?(c)", "file.c"),
		glob(true, "*.c?(c)", "file.cc"),
		glob(false, "*(@(a))a@(c)", "baaac"),
		glob(false, "*(@(a))a@(c)", "c"),
		glob(true, "*(@(a))a@(c)", "aaac"),
		glob(true, "*(@(a))a@(c)", "aac"),
		glob(true, "*(@(a))a@(c)", "ac"),

		glob(true, "*(*(f)*(o))", "fffooofoooooffoofffooofff"),
		glob(false, "*(*(of*(o)x)o)", "ofoooxoofxoofoooxoofxofo"),
		glob(true, "*(*(of*(o)x)o)", "ofoooxoofxo"),
		glob(true, "*(*(of*(o)x)o)", "ofoooxoofxoofoooxoofxo"),
		glob(true, "*(*(of*(o)x)o)", "ofoooxoofxoofoooxoofxoo"),
		glob(true, "*(*(of*(o)x)o)", "ofoooxoofxoofoooxoofxooofxofxo"),
		glob(true, "*(*(of*(o)x)o)", "ofxoofxo"),

		glob(false, "*(0|1|3|5|7|9)", "\"\""),
		glob(false, "*(0|1|3|5|7|9)", "2468"),
		glob(true, "*(0|1|3|5|7|9)", ""),
		glob(true, "*(0|1|3|5|7|9)", "137577991"),
		glob(true, "*(a)", "a"),
		glob(false, "*(a+|b)", "abef"),
		glob(false, "*(a|b\\[)", "*(a|b[)"),
		glob(false, "*(a|b\\[)", "foo"),
		glob(true, "*(b+(c)d|e*(f)g?|?(h)i@(j|k))", "egzefffgzbcdij"),

		glob(false, "*(f*(o))", "foooofofx"),
		glob(false, "*(f*(o))", "ofooofoofofooo"),
		glob(false, "*(f*(o))", "xfoooofof"),
		glob(true, "*(f*(o))", "ffo"),
		glob(true, "*(f*(o))", "fofo"),
		glob(true, "*(f*(o))", "fooofoofofooo"),
		glob(true, "*(f*(o))", "foooofo"),
		glob(true, "*(f*(o))", "foooofof"),
		glob(false, "*(f*(o)x)", "foooxfooxofoxfooox"),
		glob(true, "*(f*(o)x)", "foooxfooxfoxfooox"),
		glob(true, "*(f*(o)x)", "foooxfooxfxfooox"),
		glob(false, "*(f+(o))", "foooofof"),
		glob(true, "*(fo|foo)", "fofoofoofofoo"),
		glob(true, "*(of+(o))", "ofoofo"),
		glob(true, "*(of+(o)|f)", "ofoofo"),
		glob(true, "*(of|oof+(o))", "oofooofo"),
		glob(false, "*(oxf+(ox))", "oxfoxfox"),
		glob(true, "*(oxf+(ox))", "oxfoxoxfox"),

		glob(false, "*\\;[1-9]*([0-9])", "\"MS.FILE;\""),
		glob(false, "*\\;[1-9]*([0-9])", "\"MS.FILE;13\""),
		glob(false, "*\\;[1-9]*([0-9])", "\"MS.FILE\""),
		glob(false, "*\\;[1-9]*([0-9])", "\"VMS.FILE;\""),
		glob(true, "+(??)/..?(/)", "ab/../"),
		glob(true, "+(??|a*)/..?(/)", "ab/../"),
		glob(true, "+(?b)/..?(/)", "ab/../"),
		glob(true, "+(?b|?b)/..?(/)", "ab/../"),
		glob(false, "+()c", "abc"),
		glob(false, "+()x", "abc"),
		glob(true, "+([!/])/..?(/)", "ab/../"),
		glob(true, "+([!/])/..@(/)", "ab/../"),
		glob(true, "+([!/])/../", "ab/../"),
		glob(true, "+([^/])/..?(/)", "ab/../"),
		glob(true, "+([^/])/../", "ab/../"),
		glob(false, "+([0-7])", "09"),
		glob(true, "+([0-7])", "0377"),
		glob(true, "+([0-7])", "07"),
		glob(true, "+(*)c", "abc"),
		glob(false, "+(*)x", "abc"),
		glob(true, "+(a)", "a"),
		glob(true, "+(a*)/..?(/)", "ab/../"),
		glob(true, "+(a|b\\[)*", "abcx"),
		glob(true, "+(a|b\\[)*", "b[c"),
		glob(true, "+(ab)/..?(/)", "ab/../"),

		glob(false, "a??b", "a"),
		glob(false, "a??b", "aa"),
		glob(false, "a??b", "aab"), // micromatch extglob disagrees
		glob(false, "a?(a|b)", "bb"),
		glob(true, "a?(a|b)", "a"),
		glob(false, "a?(b*)", "ax"),
		glob(false, "a?(x)", "ab"),
		glob(false, "a?(x)", "ba"),
		glob(true, "a?(x)", "a"),
		glob(true, "a?(x)", "ax"),
		glob(false, "a[b*(foo|bar)]d", "abc"),
		glob(false, "a[b*(foo|bar)]d", "acd"),
		glob(true, "a[b*(foo|bar)]d", "abd"),
		glob(false, "a*?(x)", "ba"),
		glob(true, "a*?(x)", "a"),
		glob(true, "a*?(x)", "ab"),
		glob(true, "a*?(x)", "ax"),
		glob(false, "a\\(*b", "ab"),
		glob(true, "a\\(*b", "a((((b"),
		glob(true, "a\\(*b", "a((b"),
		glob(true, "a\\(*b", "a(b"),

		glob(false, "a+(b|c)d", "abc"),
		glob(true, "a+(b|c)d", "abd"),
		glob(true, "a+(b|c)d", "acd"),
		glob(false, "ab?*(e|f)", "123abc"),
		glob(false, "ab?*(e|f)", "ab"),
		glob(false, "ab?*(e|f)", "abcdef"),
		glob(false, "ab?*(e|f)", "abcfefg"),
		glob(false, "ab?*(e|f)", "acd"),
		glob(true, "ab?*(e|f)", "abcfef"),
		glob(true, "ab?*(e|f)", "abd"),
		glob(true, "ab?*(e|f)", "abef"),
		glob(false, "ab*(e|f)", "abcdef"),
		glob(false, "ab*(e|f)", "abcfef"),
		glob(false, "ab*(e|f)", "abcfefg"),
		glob(true, "ab*(e|f)", "ab"),
		glob(true, "ab*(e|f)", "abef"),

		glob(true, "ab**", "ab"),
		glob(true, "ab**", "abcdef"),
		glob(true, "ab**", "abcfef"),
		glob(true, "ab**", "abcfefg"),
		glob(true, "ab**", "abef"),
		glob(true, "ab**(e|f)", "ab"),
		glob(true, "ab**(e|f)", "abab"),
		glob(true, "ab**(e|f)", "abcdef"),
		glob(true, "ab**(e|f)", "abcfef"),
		glob(true, "ab**(e|f)", "abcfefg"),
		glob(true, "ab**(e|f)", "abd"),
		glob(true, "ab**(e|f)", "abef"),
		glob(false, "ab**(e|f)g", "ab"),
		glob(false, "ab**(e|f)g", "abcdef"),
		glob(false, "ab**(e|f)g", "abcfef"),
		glob(false, "ab**(e|f)g", "abef"),
		glob(true, "ab**(e|f)g", "abcfefg"),

		glob(false, "ab***ef", "ab"),
		glob(false, "ab***ef", "abcfefg"),
		glob(true, "ab***ef", "abcdef"),
		glob(true, "ab***ef", "abcfef"),
		glob(true, "ab***ef", "abef"),
		glob(false, "ab*+(e|f)", "ab"),
		glob(false, "ab*+(e|f)", "abcfefg"),
		glob(true, "ab*+(e|f)", "abcdef"),
		glob(true, "ab*+(e|f)", "abcfef"),
		glob(true, "ab*+(e|f)", "abef"),

		glob(false, "ab*d+(e|f)", "123abc"),
		glob(false, "ab*d+(e|f)", "ab"),
		glob(false, "ab*d+(e|f)", "abcfef"),
		glob(false, "ab*d+(e|f)", "abcfefg"),
		glob(false, "ab*d+(e|f)", "abd"),
		glob(false, "ab*d+(e|f)", "abef"),
		glob(false, "ab*d+(e|f)", "acd"),
		glob(true, "ab*d+(e|f)", "abcdef"),
		glob(false, "b?(a|b)", "a"),
		glob(true, "b?(a|b)", "ba"),
		glob(false, "b?*(e|f)", "ab"),
		glob(false, "b?*(e|f)", "abcdef"),
		glob(false, "b?*(e|f)", "abcfef"),
		glob(false, "b?*(e|f)", "abcfefg"),
		glob(false, "b?*(e|f)", "abef"),

		glob(false, "no-file+(a*(c)|b)stuff", "abc"),
		glob(false, "no-file+(a|b)stuff", "abc"),
		glob(false, "para?([345]|99)1", "para381"),
		glob(true, "para?([345]|99)1", "para991"),
		glob(false, "para@(chute|graph)", "paramour"),
		glob(true, "para@(chute|graph)", "paragraph"),
		glob(false, "para*([0-9])", "paragraph"),
		glob(true, "para*([0-9])", "para"),
		glob(true, "para*([0-9])", "para13829383746592"),
		glob(false, "para+([0-9])", "para"),
		glob(true, "para+([0-9])", "para987346523"),
	} {
		t.Run("Extended", func(t *testing.T) {
			g, err := Compile(test.pattern, test.delimiters...)
			if err != nil {
				t.Fatal(err)
			}
			result := g.Match(test.match)
			if result != test.should {
				t.Errorf(
					"pattern %q matching %q should be %v but got %v\n%s",
					test.pattern, test.match, test.should, result, g.r,
				)
			}
		})
	}
}

func TestNegationGlob(t *testing.T) {
	for _, test := range []test{
		// test that the dummy strings work
		glob(true, "\\\\$..!(!())", "\\$.."),

		// since the method used here is based on micromatch/extglob,
		// i've included a fix for https://github.com/micromatch/extglob/issues/10
		glob(false, "!(*.js|*.json)", "a.js"),
		glob(true, "!(*.js|*.json)", "a.js.gz"),
		glob(true, "!(*.js|*.json)", "a.json.gz"),
		glob(true, "!(*.js|*.json)", "a.gz"),
		glob(false, "!(*.js|*.json)", "a.js"),

		glob(true, "a*!(x)", "a"),
		glob(true, "a*!(x)", "ab"),
		glob(false, "a*!(x)", "ba"),
		glob(true, "a*!(x)", "ax"),
		glob(true, "a!(x)", "a"),
		glob(true, "a!(x)", "ab"),
		glob(false, "a!(x)", "ba"),
		glob(false, "a!(x)", "ax"),

		glob(true, "!(x)", "foo"),
		glob(false, "!(**/)", "foo/"),
		glob(true, "!(**/)", "foo/bar"),
		glob(true, "*(*/*)!(*/)", "foo/bar", '/'),
		glob(false, "*!(*/)", "foo/", '/'),
		glob(false, "*(*/)!(*/)", "foo/", '/'),
		glob(false, "!(*/)", "foo/"),

		glob(true, "!(x)", "foo/bar"),
		glob(false, "!(x)", "foo/bar", '/'),
		glob(true, "!(x)*", "foo"),
		glob(false, "!(foo)", "foo"),
		glob(true, "!(!(foo))", "foo"),
		glob(false, "!(!(!(foo)))", "foo"),
		glob(true, "!(!(!(!(foo))))", "foo"),
		glob(false, "!(foo)*", "foo"), // Bash 4.3 disagrees!
		glob(true, "!(foo)", "foobar"),
		glob(false, "!(foo)*", "foobar"),        // Bash 4.3 disagrees!
		glob(false, "!(*.*).!(*.*)", "moo.cow"), // Bash 4.3 disagrees!
		glob(false, "!(*.*).!(*.*)", "mad.moo.cow"),
		glob(false, "mu!(*(c))?.pa!(*(z))?", "mucca.pazza"),
		glob(true, "!(f)", "fff"),
		glob(true, "*(!(f))", "fff"),
		glob(true, "+(!(f))", "fff"),
		glob(true, "!(f)", "ooo"),
		glob(true, "*(!(f))", "ooo"),
		glob(true, "+(!(f))", "ooo"),
		glob(true, "!(f)", "foo"),
		glob(true, "*(!(f))", "foo"),
		glob(true, "+(!(f))", "foo"),
		glob(false, "!(f)", "f"),
		glob(false, "*(!(f))", "f"),
		glob(false, "+(!(f))", "f"),
		glob(true, "@(!(z*)|*x)", "foot"),
		glob(false, "@(!(z*)|*x)", "zoot"),
		glob(true, "@(!(z*)|*x)", "foox"),
		glob(true, "@(!(z*)|*x)", "zoox"),
		glob(false, "*(!(foo))", "foo"), // Bash 4.3 disagrees!
		glob(false, "!(foo)b*", "foob"),
		glob(false, "!(foo)b*", "foobb"), // Bash 4.3 disagrees!

		glob(true, "*.!(js|css)", "bar.min.js"),
		glob(false, "!*.+(js|css)", "bar.min.js"),
		glob(true, "*.+(js|css)", "bar.min.js"),

		glob(true, "*(*.json|!(*.js))", "other.bar"),
		glob(true, "*(*.json|!(*.js))*", "other.bar"),
		glob(false, "!(*(*.json|!(*.js)))*", "other.bar"),
		glob(true, "+(*.json|!(*.js))", "other.bar"),
		glob(true, "@(*.json|!(*.js))", "other.bar"),
		glob(true, "?(*.json|!(*.js))", "other.bar"),

		glob(false, "*.!(js)*.!(xy)", "asd.js.xyz"),
		glob(false, "*.!(js)*.!(xy)*", "asd.js.xyz"),
		glob(false, "*.!(js)*.!(xyz)", "asd.js.xyz"),
		glob(false, "*.!(js)*.!(xyz)*", "asd.js.xyz"),
		glob(false, "*.!(js).!(xy)", "asd.js.xyz"),
		glob(false, "*.!(js).!(xy)*", "asd.js.xyz"),
		glob(false, "*.!(js).!(xyz)", "asd.js.xyz"),
		glob(false, "*.!(js).!(xyz)*", "asd.js.xyz"),

		glob(true, "*.!(j)", "a-integration-test.js"),
		glob(false, "*.!(js)", "a-integration-test.js"),
		glob(false, "!(*-integration-test.js)", "a-integration-test.js"),
		glob(true, "*-!(integration-)test.js", "a-integration-test.js"),
		glob(false, "*-!(integration)-test.js", "a-integration-test.js"),
		glob(true, "*!(-integration)-test.js", "a-integration-test.js"),
		glob(true, "*!(-integration-)test.js", "a-integration-test.js"),
		glob(true, "*!(integration)-test.js", "a-integration-test.js"),
		glob(true, "*!(integration-test).js", "a-integration-test.js"),
		glob(true, "*-!(integration-test).js", "a-integration-test.js"),
		glob(true, "*-!(integration-test.js)", "a-integration-test.js"),
		glob(false, "*-!(integra)tion-test.js", "a-integration-test.js"),
		glob(false, "*-integr!(ation)-test.js", "a-integration-test.js"),
		glob(false, "*-integr!(ation-t)est.js", "a-integration-test.js"),
		glob(false, "*-i!(ntegration-)test.js", "a-integration-test.js"),
		glob(true, "*i!(ntegration-)test.js", "a-integration-test.js"),
		glob(true, "*te!(gration-te)st.js", "a-integration-test.js"),
		glob(false, "*-!(integration)?test.js", "a-integration-test.js"),
		glob(true, "*?!(integration)?test.js", "a-integration-test.js"),

		glob(true, "*!(js)", "foo.js.js"),
		glob(true, "*!(.js)", "foo.js.js"),
		glob(true, "*!(.js.js)", "foo.js.js"),
		glob(true, "*!(.js.js)*", "foo.js.js"),
		glob(false, "*(.js.js)", "foo.js.js"),
		glob(true, "**(.js.js)", "foo.js.js"),
		glob(true, "*(!(.js.js))", "foo.js.js"),
		glob(false, "*.!(js)*.!(js)", "foo.js.js"),
		glob(false, "*.!(js)+", "foo.js.js"),
		glob(true, "!(*(.js.js))", "foo.js.js"),
		glob(true, "*.!(js)", "foo.js.js"),
		glob(false, "*.!(js)*", "foo.js.js"),    // Bash 4.3 disagrees,
		glob(false, "*.!(js)*.js", "foo.js.js"), // Bash 4.3 disagrees,

		glob(true, "*/**(.*)", "a/foo.js.js"),
		glob(true, "*/**(.*.*)", "a/foo.js.js"),
		glob(true, "a/**(.*.*)", "a/foo.js.js"),
		glob(true, "*/**(.js.js)", "a/foo.js.js"),
		glob(true, "a/f*(!(.js.js))", "a/foo.js.js"),
		glob(true, "a/!(*(.*))", "a/foo.js.js"),
		glob(true, "a/!(+(.*))", "a/foo.js.js"),
		glob(true, "a/!(*(.*.*))", "a/foo.js.js"),
		glob(true, "*/!(*(.*.*))", "a/foo.js.js"),
		glob(true, "a/!(*(.js.js))", "a/foo.js.js"),

		glob(true, "*(*.json|!(*.js))", "testjson.json"),
		glob(true, "+(*.json|!(*.js))", "testjson.json"),
		glob(true, "@(*.json|!(*.js))", "testjson.json"),
		glob(true, "?(*.json|!(*.js))", "testjson.json"),

		glob(false, "*(*.json|!(*.js))", "foojs.js"), // Bash 4.3 disagrees
		glob(true, "*(*.json|!(*.js))*", "foojs.js"),
		glob(false, "+(*.json|!(*.js))", "foojs.js"), // Bash 4.3 disagrees
		glob(false, "@(*.json|!(*.js))", "foojs.js"),
		glob(false, "?(*.json|!(*.js))", "foojs.js"),

		glob(true, "!(*.a|*.b|*.c)", "a"),
		glob(false, "!(*.[a-b]*)", "a.a"),
		glob(false, "!(*.a|*.b|*.c)", "a.a"),
		glob(false, "!(*[a-b].[a-b]*)", "a.a"),
		glob(false, "!*.(a|b)", "a.a"),
		glob(false, "!*.(a|b)*", "a.a"),
		glob(false, "*.!(a)", "a.a"),
		glob(false, "*.+(b|d)", "a.a"),
		glob(false, "!(*.[a-b]*)", "a.a.a"),
		glob(false, "!(*[a-b].[a-b]*)", "a.a.a"),
		glob(false, "!*.(a|b)", "a.a.a"),
		glob(false, "!*.(a|b)*", "a.a.a"),
		glob(true, "!(*.a|*.b|*.c)", "a.abcd"), // micromatch/extglob disagrees
		glob(true, "!(*.a|*.b|*.c)", "c.cbad"), // but interestingly, agrees here
		glob(false, "!(*.a|*.b|*.c)*", "a.abcd"),
		glob(true, "*.!(a|b|c)", "a.abcd"), // micromatch/extglob disagrees
		glob(false, "*.!(a|b|c)*", "a.abcd"),
		glob(false, "!(*.*)", "a.b"),
		glob(false, "!(*.[a-b]*)", "a.b"),
		glob(false, "!(*[a-b].[a-b]*)", "a.b"),
		glob(false, "!*.(a|b)", "a.b"),
		glob(false, "!*.(a|b)*", "a.b"),
		glob(false, "!(*.[a-b]*)", "a.bb"),
		glob(false, "!(*[a-b].[a-b]*)", "a.bb"),
		glob(false, "!*.(a|b)", "a.bb"),
		glob(false, "!*.(a|b)*", "a.bb"),
		glob(false, "!*.(a|b)", "a.ccc"),
		glob(false, "!*.(a|b)*", "a.ccc"),
		glob(false, "*.+(b|d)", "a.ccc"),
		glob(false, "!(*.js)", "a.js"),
		glob(false, "*.!(js)", "a.js"),
		glob(false, "!(*.js)", "a.js.js"),
	} {
		t.Run("Negated", func(t *testing.T) {
			g, err := Compile(test.pattern, test.delimiters...)
			if err != nil {
				t.Fatal(err)
			}
			result := g.Match(test.match)
			if result != test.should {
				t.Errorf(
					"pattern %q matching %q should be %v but got %v\n%s",
					test.pattern, test.match, test.should, result, g.r,
				)
			}
		})
	}
}
