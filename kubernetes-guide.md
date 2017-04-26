I. Kubernetes basic
1. Configuration & install kubernetes on Local machine (Ubuntu)
* Create a Minikube cluster. Minikube is a tool to run Kubernetes locally. Minibuke create a single node Kubernetes cluster in VM on Local machine like using Virtualbox
```bash
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.15.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
```
Run ```minikube version``` to check whether it was installed successfully.
* Intall Kubectl
```bash
wget https://storage.googleapis.com/kubernetes-release/release/v1.4.4/bin/linux/amd64/kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin/kubectl
```
* Ok. We are done with installing of Minikube and kubectl. So first of all, start minikube and config kubectl to be able to interact with minikube (cluster) by these following commands. 
```minikube start``` - this will create a virtual machine by using virtual box. 
```kubectl config use-context minikube``` - to set the Minikube context. The context is what determines which cluster kubectl is interacting with
```kubectl cluster-info``` - to verify that kubectl is configured to communicate with the cluster.
2. Kubectl commands
* To view nodes in cluster
```kubectl get nodes```
* To create a deployment
```bash
kubectl run deployment_name --image=docker_image --port=port
```
* To see list of deployments
```kubectl get deployments```
* To view list of pods
```kubectl get pods```
* To view what containers are inside those pods and what images are used to build those containers
```kubectl describe pods``` or ```kubectl describe pod specific_pod_name```
* To view the container logs
```kubectl logs pod_name```
* To excute command on the container
```kubectl exec -ti pod_name /bin/bash``` - you can get inside container by this.
* To create a service so that the app can be accessed from outside of cluster
```bash
kubectl expose deployment/deployment_name --type=type --port port
```
type: could be NodePort, LoadBalancer
port: is a port of container in pod.
* To see list services - we can see the service received an unique cluster-IP, an internal port and external-IP (an IP of the Node)
```kubectl get services```
* To find out what port was opened externally (by the NodePort option) we'll run
```
kubectl describe services/service_name
```
* Using ```curl``` to test the app is exposed outside of the cluster
```curl ip_of_node:node_port```
* Using labels
** To see the name of label of pod
```kubectl describe deployment```
** To query list of pods using label
```kubectl get pods -l run=label_name```
** To query list of services using label
```kubectl get services -l run=label_name```
** To update label of pod
```kubectl label pod pod_name new_label_name```
** To delete a service
```kubectl delete service -l run=label_name```
** To scale a deployment
*** To scale a deployment
```kubectl scale deployments/deployment_name --replicas=number_instances_wanted_to_scale_to```
*** To see number of pods after scaling
```kubectl get pods -o wide```
*** To check the service is load-balancing the traffic, we use ```curl ip_of_node:port_of_node```. This will hit a difference pod with every request.
*** To scale down
```bash
kubectl scale deployments/kubernetes-bootcamp --replicas=number_instances_wanted_to_scale_to
```
Check list of pods after scaling down ```kubectl get pods -o wide```. this confirms that 2 pods were terminated
* To update the version of the app
  * Update the image of the application to another version
  ```
  kubectl set image deployments/deployment_name=docker_image:newversion
  ```
  * Check the status of a rolling update, we will see the old one is terminating
  ```kubectl get pods```
  * To verify an update, we use ```curl ip_node:port_node```
  * To see the version of current docker image
  ```kubectl describe pods```
  * Rollback an update ( this is for when update got error, we can rollback to previous version)
  ```kubectl rollout undo deployments/deployment_name```

  3. Let's make an example to get the hang of basic kubernetes: creating a simple node.js application using docker container. Located all needed files in folder ```hellonode```
    * Create a ```server.js``` file

    ```js
    var http = require('http');
    var handleRequest = function(request, response) {
      console.log('Received request for URL: ' + request.url);
      response.writeHead(200);
      response.end('Hello World!');
    };
    var www = http.createServer(handleRequest);
    www.listen(8080);
    ```
    * You can run it locally on your machine to see it works
    ```node server.js```
    ```curl http://localhost:8080```
    The message is supposed to be "Hello world"
    * Now let's dive into it, create a Docker container image by creating a ```Dockerfile```

    ```
    FROM node:6.9.2
    EXPOSE 8080
    COPY server.js .
    CMD node server.js
    ```
    * To build the image using the same Docker host as the Minikube VM
        ```
        eval $(minikube docker-env)
        ```
    * Build the docker image using Minikube daemon
        ```
        docker build -t hello-node:v1 .
        ```
    * Create a Deployment
        A Kubernetes Deployment checks on the health of Pods and restarts the Pod's container if it terminates. Deployments are the recommended way to manage the creation and scaling of Pods. 
        Use ```kubectl run``` command to create a Deployment that manages a Pod. The Pod runs a Container based on your hello-node:v1 Docker image. 
        ```
        kubectl run hello-node --image=hello-node:v1 --port=8080
        ```
        To view the deployment:
        ```
        kubectl get deployments
        ```
        To view the pods
        ```kubectl get pods```
    * Create a service to access Pods from outside of Kubernetes cluster
        ```
        kubectl expose deployment hello-node --type=LoadBalancer
        ```
    * Run ```kubectl describe pods``` to see the node_ip
    * Run ```kubectl describe service hello-node``` to see node_port
    * We can check by using ```curl```
        ```curl node_ip:node_port```
        It should return ```hello world``` on terminal

II. Kubernetes with Yaml, json file

  Instead of creating deployment, pod, ReplicationController or Service by typing those commands above, we can put all configuration in yaml file or json file. Let's get started.
  1. Create Pod

    The sample file will be like this

    ```
    apiVersion: v1
    kind: Pod
    metadata:
      name: set_name_for_pod
      labels:
        key1: value1
        key2: value2
    spec:
      containers:
      - image: docker_image_name
        name: set_name_for_container
        env:
        - name: ...
          value: ...
        ports:
        - containerPort: port_number
          name: name_of_port
        volumeMounts:
        - name: set_name_for_volume
          mountPath: the_path_in_container
      volumes:
      - name: set_name_for_volume
        persistentVolumeClaim:
          claimName: name_of_claim
    ```

    ex: file ```mysql_pod.yaml```

    ```
    apiVersion: v1
    kind: Pod
    metadata:
      name: mysql
      labels:
        type: local
        tier: backend
    spec:
      containers:
      - image: mysql
        name: mysql
        env:
        - name: WORDPRESS_ROOT_PASSWORD
          valueFrom:
              secretKeyRef:
                name: mysql-pass
                key: password.txt # password was saved in a file
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-persistent-storage
            mountPath: /var/lib/mysql
      volumes:
      - name: mysql-persistent-storage
        persistentVolumeClaim:
          claimName: mysql-pv-claim # this claim has created already
    ```

    To create pod by run this command
    ```kubectl create -f mysql_pod.yaml```

    Learn more here [Pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/)
  2. Create Deployment

    A Deployment provides declarative updates for Pods or Replica Sets (the next generation ReplicationController). A typical use case is:
    - Create a Deployment to bring up a Replica Set or Pods, 
    - Check the status of a Deployment to see if it succeeds or not
    - when to need to use a new image, update Deployment to recreate the Pods
    - Rollback to an earlier Deployment revision if the current is not stable. 
    - Pause and resume a Deployment
    
    Sample yaml file:

    ```
    apiVersion: v1
    kind: Deployment
    metadata:
      name: set_deployment_name
    spec:
      replicas: number_of_pods
      # this creates pods so the code is the same as when creating pod above
      template:
        metadata:
          labels:
            key1: value1
            key2: value2
        spec:
          containers:
          - name: container
          .... # same as creating pod above
    ```

    ex: file ```mysql-deployment.yaml```

    ```
    apiVersion: v1
    kind: Deployment
    metadata: 
      name: mysql-deployment
    spec:
      replicas: 2
      template:
        name: mysql
        labels:
          type: local
          tier: backend
      spec:
        containers:
        - image: mysql
          name: mysql
          env:
          - name: WORDPRESS_ROOT_PASSWORD
            valueFrom:
                secretKeyRef:
                  name: mysql-pass
                  key: password.txt # password was saved in a file
          ports:
          - containerPort: 3306
            name: mysql
          volumeMounts:
          - name: mysql-persistent-storage
              mountPath: /var/lib/mysql
        volumes:
        - name: mysql-persistent-storage
          persistentVolumeClaim:
            claimName: mysql-pv-claim # this claim has created already
    ```

    Run this command to create deployment
    ```kubectl create -f mysql-deployment.yaml```

    To get pods, run ```kubectl get pods```, it should show 2 mysql pods.

    To update a Deployment, run
    ```kubectl set image deployment/mysql-deployment mysql=mysql:5.6```

    To rollback to a previous revision
    ```kubectl rollout undo deployment/mysql-deployment```

    Learn more here [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
  
  3. ReplicationController

    A Replication Controller makes sure that a pod or homogeneous set of pods are always up and available, the pods maintained by a ReplicationController are automatically replaced if they fail, get deleted or are terminated. 

    Sample yaml file

    ```
    apiVersion: v1
    kind: ReplicationController
    metadata:
      name: name_of_rc
    spec:
      replicas: number_of_pods
      # this will select the pods which have label ```app=name_of_app```
      # to be maintained by this ReplicationController
      selector: 
        app: name_of_app
      template: # this creates pods so the code is the same as when creating pod
        name: ...
        labels:
          key1: value1
      spec:
        containers:
        - name: container_name
        ... # same as creating pod
    ```

    ex: ```replication.yaml```

    ```
    apiVersion: v1
    kind: ReplicationController
    metadata:
      name: nginx
    spec:
      replicas: 3
      selector:
        app: nginx
      template:
        metadata:
          name: nginx
          labels:
            app: nginx
        spec:
          containers:
          - name: nginx
            image: nginx
            ports:
            - containerPort: 80
    ```

    Run ```kubectl create -f ./replication.yaml``` to create replication controller

    To see replication run ```kubectl get rc```

    To check status of replication controller run
    ```kubectl describe rc nginx```

    Learn more here [Replication Controller](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/)

  4. Service
    A Service in Kubernetes is a REST object, similar to a Pod. Like all of the REST objects, a Service definition can be POSTed to the apiserver to create a new instance. For example, suppose you have a set of Pods that each expose port ```9376``` and carry a label ```"app=MyApp"```.

    Sample yaml file

    ```
    kind: Service
    apiVersion: v1
    metadata:
      name: my-service
    spec:
      selector:
        app: MyApp
      ports:
        - protocol: TCP
          port: 80
          targetPort: 9376
      type: LoadBalancer
    ```

    This specification will create a new ```Service``` named ```my-service``` which targets port ```9376``` on any Pod with the ```app=MyApp``` label. Type ```loadBalancer``` is to create an external load balancer. Type NodePort
    Kubernetes Services support TCP and UDP for protocols. The default is TCP.

    Learn more here [Service](https://kubernetes.io/docs/concepts/services-networking/service/)

  5. Persistent Volumes

    5.1. Create a PersistentVolume by yaml file

      A ```PersistentVolume``` (PV) is a piece of networked storage in the cluster that has been provisioned by an administrator. It is a resource in the cluster just like a node is a cluster resource. 

      ```
      kind: PersistentVolume
      apiVersion: v1
      metadata:
        name: persistentVolumeName
        labels:
          name: name_of_label (it can be same name as persistentVolumeName)
      spec:
        capacity:
          storage: set_a_specific_size_of_volume (example 10Gi)
        accessModes:
          - ReadWriteOnce
        claimRef:
          namespace: name_of_namespace_of_persistentVolumeClaim
          name: name_of_persistentVolumeClaim
        hostPath:
          path: "folder_path_on_node"
      ```

      Explain:
        - Capacity: a PV will have a specific storage capacity. This is set using the PV's ```capacity``` attribute. 
        - accessModes: a ```PersistentVolume``` can be mounted on a host in any way supported by the resource provider. The access modes are:
          * ReadWriteOnce: the volume can be mounted as read-write by a single node
          * ReadOnlyMany: the volume can be mounted read-only by many nodes
          * ReadWriteMany: the volume can be mounted as read-write by many nodes
        *Important!* A volume can only be mounted using one access mode at a time, even if it supports many.

    5.2. Create a PersistentVolumeClaim by yaml file.

      A ```PersistentVolumeClaim``` (PVC) is a request for storage by a user. It is similar to a pod. Pods consume node resources and PVCs consume PV resource. Pods can request specific levels of resources (CPU and Memory). Claims can request specific size and access modes 

      ```
      kind: PersistentVolumeClaim
      apiVersion: v1
      metadata: 
        name: myClaim
      spec:
        accessMode:
          - ReadWriteOnce
        resources:
          requests: 
            storage: 8Gi
      ```

      Explain:
        - Access Modes: Claims use the same conventions as volumes when requesting storage with specific access modes
        - Resources: Claims, like pods, can request specific quantities of a resource. In this case, the request is for storage. The same resource model applies to both volumes and claims

    5.3. Claims as Volumes:

      Pods access storage by using the claim as a volume. Claims must exist in the same namespace as the pod using the claim. The cluster finds the claim in the pod's namespace and uses it to get the ```PersistentVolume``` backing the claim. The volume is then mounted to the host and into the pod.

      ```
      kind: Pod
      apiVersion: v1
      metadata:
        name: podName
      spec:
        containers:
          name: containerName
          image: dockerImage
          volumeMounts:
          - mountPath: "path_in_container"
            name: mypd
        volumes:
          name: mypd
          persistentVolumeClaim:
            claimName: myclaim
      ```

      Remember name of volumeMounts in container should be same name as name of volumes

      Learn more here [Persistent Volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
  
