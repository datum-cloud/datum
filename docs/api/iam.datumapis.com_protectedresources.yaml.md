# API Reference

Packages:

- [iam.datumapis.com/v1alpha1](#iamdatumapiscomv1alpha1)

# iam.datumapis.com/v1alpha1

Resource Types:

- [ProtectedResource](#protectedresource)




## ProtectedResource
<sup><sup>[↩ Parent](#iamdatumapiscomv1alpha1 )</sup></sup>






ProtectedResource is the Schema for the protectedresources API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>iam.datumapis.com/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>ProtectedResource</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#protectedresourcespec">spec</a></b></td>
        <td>object</td>
        <td>
          ProtectedResourceSpec defines the desired state of ProtectedResource<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#protectedresourcestatus">status</a></b></td>
        <td>object</td>
        <td>
          ProtectedResourceStatus defines the observed state of ProtectedResource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProtectedResource.spec
<sup><sup>[↩ Parent](#protectedresource)</sup></sup>



ProtectedResourceSpec defines the desired state of ProtectedResource

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          The kind of the resource.
This will be in the format `Workload`.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>permissions</b></td>
        <td>[]string</td>
        <td>
          A list of permissions that are associated with the resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>plural</b></td>
        <td>string</td>
        <td>
          The plural form for the resource type, e.g. 'workloads'. Must follow
camelCase format.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#protectedresourcespecserviceref">serviceRef</a></b></td>
        <td>object</td>
        <td>
          ServiceRef references the service definition this protected resource belongs to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>singular</b></td>
        <td>string</td>
        <td>
          The singular form for the resource type, e.g. 'workload'. Must follow
camelCase format.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#protectedresourcespecparentresourcesindex">parentResources</a></b></td>
        <td>[]object</td>
        <td>
          A list of resources that are registered with the platform that may be a
parent to the resource. Permissions may be bound to a parent resource so
they can be inherited down the resource hierarchy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProtectedResource.spec.serviceRef
<sup><sup>[↩ Parent](#protectedresourcespec)</sup></sup>



ServiceRef references the service definition this protected resource belongs to.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the resource name of the service definition.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ProtectedResource.spec.parentResources[index]
<sup><sup>[↩ Parent](#protectedresourcespec)</sup></sup>



ParentResourceRef defines the reference to a parent resource

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          APIGroup is the group for the resource being referenced.
If APIGroup is not specified, the specified Kind must be in the core API group.
For any other third-party types, APIGroup is required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProtectedResource.status
<sup><sup>[↩ Parent](#protectedresource)</sup></sup>



ProtectedResourceStatus defines the observed state of ProtectedResource

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#protectedresourcestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions provide conditions that represent the current status of the ProtectedResource.<br/>
          <br/>
            <i>Default</i>: [map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for control plane to reconcile reason:Unknown status:Unknown type:Ready]]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the most recent generation observed for this ProtectedResource. It corresponds to the
ProtectedResource's generation, which is updated on mutation by the API Server.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ProtectedResource.status.conditions[index]
<sup><sup>[↩ Parent](#protectedresourcestatus)</sup></sup>



Condition contains details for one aspect of the current state of this API Resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
