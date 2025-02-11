# Azure

You can deploy Pachyderm in a new or existing Microsoft® Azure® Kubernetes
Service environment and use Azure's resource to run your Pachyderm
workloads. 
To deploy Pachyderm to AKS, you need to:

1. [Install Prerequisites](#install-prerequisites)
2. [Deploy Kubernetes](#deploy-kubernetes)
3. [Deploy Pachyderm](#deploy-pachyderm)

## Install Prerequisites

Before you can deploy Pachyderm on Azure, you need to configure a few
prerequisites on your client machine. If not explicitly specified, use the
latest available version of the components listed below.
Install the following prerequisites:

* [Azure CLI 2.0.1 or later](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
* [jq](https://stedolan.github.io/jq/download/)
* [kubectl](https://docs.microsoft.com/cli/azure/aks?view=azure-cli-latest#az_aks_install_cli)
* [pachctl](#install-pachctl)

### Install `pachctl`

 `pachctl` is a primary command-line utility for interacting with Pachyderm clusters.
 You can run the tool on Linux®, macOS®, and Microsoft® Windows® 10 or later operating
 systems and install it by using your favorite command line package manager.
 This section describes how you can install `pachctl` by using
 `brew` and `curl`.

 If you are installing `pachctl` on Windows, you need to first install
 Windows Subsystem (WSL) for Linux.

 To install `pachctl`, complete the following steps:

 * To install on macOS by using `brew`, run the following command:

   ```bash
   $ brew tap pachyderm/tap && brew install pachyderm/tap/pachctl@1.9
   ```
 * To install on Linux 64-bit or Windows 10 or later, run the following command:

   ```bash
   $ curl -o /tmp/pachctl.deb -L https://github.com/pachyderm/pachyderm/releases/download/v1.9.3/pachctl_1.9.3_amd64.deb &&  sudo dpkg -i /tmp/pachctl.deb
   ```

 1. Verify your installation by running `pachctl version`:

    ```bash
    $ pachctl version --client-only
    COMPONENT           VERSION
    pachctl             1.9.0
    ```

## Deploy Kubernetes

You can deploy Kubernetes on Azure by following the official [Azure Container Service documentation](https://docs.microsoft.com/azure/aks/tutorial-kubernetes-deploy-cluster) or by
following the steps in this section. When you deploy Kubernetes on Azure,
you need to specify the following parameters:

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#ccc;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#ccc;color:#333;background-color:#fff;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#ccc;color:#333;background-color:#f0f0f0;}
.tg .tg-0pky{border-color:inherit;text-align:left;vertical-align:top}
</style>
<table class="tg">
  <tr>
    <th class="tg-0pky">Variable</th>
    <th class="tg-0pky">Description</th>
  </tr>
  <tr>
    <td class="tg-0pky">RESOURCE_GROUP</td>
    <td class="tg-0pky">A unique name for the resource group where Pachyderm is deployed. For example, `pach-resource-group`.</td>
  </tr>
  <tr>
    <td class="tg-0pky">LOCATION</td>
    <td class="tg-0pky">An Azure availability zone where AKS is available. For example, `centralus`.</td>
  </tr>
  <tr>
    <td class="tg-0pky">NODE_SIZE</td>
    <td class="tg-0pky">The size of the Kubernetes virtual machine (VM) instances. To avoid performance issues, Pachyderm recommends that you
    set this value to at least `Standard_DS4_v2` which gives you 8 CPUs, 28 Gib of Memory, 56 Gib SSD.</td>
  </tr>
  <tr>
    <td class="tg-0pky">CLUSTER_NAME</td>
    <td class="tg-0pky">A unique name for the Pachyderm cluster. For example, `pach-aks-cluster`.</td>
  </tr>
</table>

To deploy Kubernetes on Azure, complete the following steps:

1. Log in to Azure:

   ```bash
   $ az login
   Note, we have launched a browser for you to login. For old experience with
   device code, use "az login --use-device-code"
   ```

   If you have not already logged in this command opens a browser window. Log in with your Azure credentials.
   After you log in, the following message appears in the command prompt:

   ```bash
   You have logged in. Now let us find all the subscriptions to which you have access...
   [
     {
       "cloudName": "AzureCloud",
       "id": "your_id",
       "isDefault": true,
       "name": "Microsoft Azure Sponsorship",
       "state": "Enabled",
       "tenantId": "your_tenant_id",
       "user": {
         "name": "your_contact_id",
         "type": "user"
       }
     }
   ]
   ```

1. Create an Azure resource group.

   ```bash
   $ az group create --name=${RESOURCE_GROUP} --location=${LOCATION}
   ```

   **Example:**

   ```bash
   $ az group create --name="test-group" --location=centralus
   {
     "id": "/subscriptions/6c9f2e1e-0eba-4421-b4cc-172f959ee110/resourceGroups/pach-resource-group",
     "location": "centralus",
     "managedBy": null,
     "name": "pach-resource-group",
     "properties": {
       "provisioningState": "Succeeded"
     },
     "tags": null,
     "type": null
   }
   ```

1. Create an AKS cluster:

   ```bash
   $ az aks create --resource-group ${RESOURCE_GROUP} --name ${CLUSTER_NAME} --generate-ssh-keys --node-vm-size ${NODE_SIZE}
   ```

   **Example:**

   ```bash
   $ az aks create --resource-group test-group --name test-cluster --generate-ssh-keys --node-vm-size Standard_DS4_v2
   {
     "aadProfile": null,
     "addonProfiles": null,
     "agentPoolProfiles": [
       {
         "availabilityZones": null,
         "count": 3,
         "enableAutoScaling": null,
         "maxCount": null,
         "maxPods": 110,
         "minCount": null,
         "name": "nodepool1",
         "orchestratorVersion": "1.12.8",
         "osDiskSizeGb": 100,
         "osType": "Linux",
         "provisioningState": "Succeeded",
         "type": "AvailabilitySet",
         "vmSize": "Standard_DS4_v2",
         "vnetSubnetId": null
       }
     ],
   ...
   ```

1. Confirm the version of the Kubernetes server:

   ```bash
   $ kubectl version
   Client Version: version.Info{Major:"1", Minor:"13", GitVersion:"v1.13.4", GitCommit:"c27b913fddd1a6c480c229191a087698aa92f0b1", GitTreeState:"clean", BuildDate:"2019-03-01T23:36:43Z", GoVersion:"go1.12", Compiler:"gc", Platform:"darwin/amd64"}
   Server Version: version.Info{Major:"1", Minor:"13", GitVersion:"v1.13.4", GitCommit:"c27b913fddd1a6c480c229191a087698aa92f0b1", GitTreeState:"clean", BuildDate:"2019-02-28T13:30:26Z", GoVersion:"go1.11.5", Compiler:"gc", Platform:"linux/amd64"}
   ```

**See also:**

- [Azure Virtual Machine sizes](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/sizes-general)


## Add storage resources

Pachyderm requires you to deploy an object store and a persistent
volume in your cloud environment to function correctly. For best
results, you need to use faster disk drives, such as *Premium SSD
Managed Disks* that are available with the Azure Premium Storage offering.

You need to specify the following parameters when you create storage
resources:

<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;border-color:#ccc;}
.tg td{font-family:Arial, sans-serif;font-size:14px;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#ccc;color:#333;background-color:#fff;}
.tg th{font-family:Arial, sans-serif;font-size:14px;font-weight:normal;padding:10px 5px;border-style:solid;border-width:1px;overflow:hidden;word-break:normal;border-color:#ccc;color:#333;background-color:#f0f0f0;}
.tg .tg-0pky{border-color:inherit;text-align:left;vertical-align:top}
</style>
<table class="tg">
  <tr>
    <th class="tg-0pky">Variable</th>
    <th class="tg-0pky">Description</th>
  </tr>
  <tr>
    <td class="tg-0pky">STORAGE_ACCOUNT</td>
    <td class="tg-0pky">The name of the storage account where you store your data, unique in the Azure location</td>
  </tr>
  <tr>
    <td class="tg-0pky">CONTAINER_NAME</td>
    <td class="tg-0pky">The name of the Azure blob container where you store your data</td>
  </tr>
  <tr>
    <td class="tg-0pky">STORAGE_SIZE</td>
    <td class="tg-0pky">The size of the persistent volume to create in GBs. Allocate at least 10 GB.</td>
  </tr>
</table>

To create these resources, follow these steps:

1. Clone the [Pachyderm GitHub repo](https://github.com/pachyderm/pachyderm).
1. Change the directory to the root directory of the `pachyderm` repository.
1. Create an Azure storage account:

   ```bash
   $ az storage account create \
     --resource-group="${RESOURCE_GROUP}" \
     --location="${LOCATION}" \
     --sku=Premium_LRS \
     --name="${STORAGE_ACCOUNT}" \
     --kind=BlockBlobStorage
   ```
   **System response:**

   ```
   {
     "accessTier": null,
     "creationTime": "2019-06-20T16:05:55.616832+00:00",
     "customDomain": null,
     "enableAzureFilesAadIntegration": null,
     "enableHttpsTrafficOnly": false,
     "encryption": {
       "keySource": "Microsoft.Storage",
       "keyVaultProperties": null,
       "services": {
         "blob": {
           "enabled": true,
     ...
   ```

   Make sure that you set Stock Keeping Unit (SKU) to `Premium_LRS`
   and the `kind` parameter is set to `BlockBlobStorage`. This
   configuration results in a storage that uses SSDs rather than
   standard Hard Disk Drives (HDD).
   If you set this parameter to an HDD-based storage option, your Pachyderm
   cluster will be too slow and might malfunction.

1. Verify that your storage account has been successfully created:

   ```bash
   $ az storage account list
   ```

1. Build a Microsoft tool for creating Azure VMs from an image:

   ```bash
   $ STORAGE_KEY="$(az storage account keys list \
                 --account-name="${STORAGE_ACCOUNT}" \
                 --resource-group="${RESOURCE_GROUP}" \
                 --output=json \
                 | jq '.[0].value' -r
              )"
   ```

1. Find the generated key in the **Storage accounts > Access keys**
   section in the Azure Portal or by running the following command:

   ```bash
   $ az storage account keys list --account-name=${STORAGE_ACCOUNT}
   [
     {
       "keyName": "key1",
       "permissions": "Full",
       "value": ""
     }
   ]
   ```

**See Also**

- [Azure Storage](https://azure.microsoft.com/documentation/articles/storage-introduction/)


## Deploy Pachyderm

After you complete all the sections above, you can deploy Pachyderm
on Azure. If you have previously tried to run Pachyderm locally,
make sure that you are using the right Kubernetes context. Otherwise,
you might accidentally deploy your cluster on Minikube.

1. Verify cluster context:

   ```bash
   $ kubectl config current-context
   ```

   This command should return the name of your Kubernetes cluster that
   runs on Azure.

   * If you have a different contents displayed, configure `kubectl`
   to use your Azure configuration:

   ```bash
   $ az aks get-credentials --resource-group ${RESOURCE_GROUP} --name ${CLUSTER_NAME}
   Merged "${CLUSTER_NAME}" as current context in /Users/test-user/.kube/config
   ```

1. Run the following command:

   ```bash
   $ pachctl deploy microsoft ${CONTAINER_NAME} ${STORAGE_ACCOUNT} ${STORAGE_KEY} ${STORAGE_SIZE} --dynamic-etcd-nodes 1
   ```
   **Example:**

   ```bash
   $ pachctl deploy microsoft test-container teststorage <key> 10 --dynamic-etcd-nodes 1
   serviceaccount/pachyderm configured
   clusterrole.rbac.authorization.k8s.io/pachyderm configured
   clusterrolebinding.rbac.authorization.k8s.io/pachyderm configured
   service/etcd-headless created
   statefulset.apps/etcd created
   service/etcd configured
   service/pachd configured
   deployment.apps/pachd configured
   service/dash configured
   deployment.apps/dash configured
   secret/pachyderm-storage-secret configured

   Pachyderm is launching. Check its status with "kubectl get all"
   Once launched, access the dashboard by running "pachctl port-forward"
   ```

   Because Pachyderm pulls containers from DockerHub, it might take some time
   before the `pachd` pods start. You can check the status of the
   deployment by periodically running `kubectl get all`.

1. When pachyderm is up and running, get the information about the pods:

   ```sh
   $ kubectl get pods
   NAME                      READY     STATUS    RESTARTS   AGE
   dash-482120938-vdlg9      2/2       Running   0          54m
   etcd-0                    1/1       Running   0          54m
   pachd-1971105989-mjn61    1/1       Running   0          54m
   ```

   **Note:** Sometimes Kubernetes tries to start `pachd` nodes before
   the `etcd` nodes are ready which might result in the `pachd` nodes
   restarting. You can safely ignore those restarts.

1. To connect to the cluster from your local machine, such as your laptop,
set up port forwarding to enable `pachctl` and cluster communication:

   ```bash
   $ pachctl port-forward
   ```

1. Verify that the cluster is up and running:

   ```sh
   $ pachctl version
   COMPONENT           VERSION
   pachctl             1.9.0
   pachd               1.9.0
   ```
