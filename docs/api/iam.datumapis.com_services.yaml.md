# API Reference

Packages:

- [iam.datumapis.com/v1alpha1](#iamdatumapiscomv1alpha1)

# iam.datumapis.com/v1alpha1

Resource Types:

- [Service](#service)




## Service
<sup><sup>[↩ Parent](#iamdatumapiscomv1alpha1 )</sup></sup>






Service is the Schema for the services API

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
      <td>Service</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#servicespec">spec</a></b></td>
        <td>object</td>
        <td>
          ServiceSpec defines the desired state of Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#servicestatus">status</a></b></td>
        <td>object</td>
        <td>
          ServiceStatus defines the observed state of Service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Service.spec
<sup><sup>[↩ Parent](#service)</sup></sup>



ServiceSpec defines the desired state of Service

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
        <td><b><a href="#servicespecresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          List of resources offered by a service.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Service.spec.resources[index]
<sup><sup>[↩ Parent](#servicespec)</sup></sup>



ServiceResource is an entity offered by services to provide functionality to service
consumers. Resources can have actions registered that result in permissions
being created.

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
        <td><b>parentResources</b></td>
        <td>[]string</td>
        <td>
          A list of resources that are registered with the platform that may be a
parent to the resource. Permissions may be bound to a parent resource so
they can be inherited down the resource hierarchy. The resource must use
the fully qualified resource name (e.g. compute.datumapis.com/Workload).<br/>
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
        <td><b>resourceNamePatterns</b></td>
        <td>[]string</td>
        <td>
          A list of resource name patterns that may be present for the resource.<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The fully qualified name of the resource.
This will be in the format `compute.datumapis.com/Workload`.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Service.status
<sup><sup>[↩ Parent](#service)</sup></sup>



ServiceStatus defines the observed state of Service

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
        <td><b><a href="#servicestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions provide conditions that represent the current status of the Service.<br/>
          <br/>
            <i>Default</i>: [{"type": "Ready", "status": "Unknown", "reason": "Unknown", "message": "Waiting for control plane to reconcile"}]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Service.status.conditions[index]
<sup><sup>[↩ Parent](#servicestatus)</sup></sup>



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
