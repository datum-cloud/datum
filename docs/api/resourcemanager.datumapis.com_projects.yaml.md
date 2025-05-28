# API Reference

Packages:

- [resourcemanager.datumapis.com/v1alpha1](#resourcemanagerdatumapiscomv1alpha1)

# resourcemanager.datumapis.com/v1alpha1

Resource Types:

- [Project](#project)




## Project
<sup><sup>[↩ Parent](#resourcemanagerdatumapiscomv1alpha1 )</sup></sup>






Project is the Schema for the projects API.

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
      <td>resourcemanager.datumapis.com/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Project</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#projectspec">spec</a></b></td>
        <td>object</td>
        <td>
          ProjectSpec defines the desired state of Project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#projectstatus">status</a></b></td>
        <td>object</td>
        <td>
          ProjectStatus defines the observed state of Project.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Project.spec
<sup><sup>[↩ Parent](#project)</sup></sup>



ProjectSpec defines the desired state of Project.

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
        <td><b><a href="#projectspecparent">parent</a></b></td>
        <td>object</td>
        <td>
          A reference to the project's parent in the resource hierarchy.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Project.spec.parent
<sup><sup>[↩ Parent](#projectspec)</sup></sup>



A reference to the project's parent in the resource hierarchy.

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
        <td><b><a href="#projectspecparentresourceref">resourceRef</a></b></td>
        <td>object</td>
        <td>
          Resource is a reference to the parent of the project. Must be a valid
resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>external</b></td>
        <td>string</td>
        <td>
          External is a reference to the parent of the project. Must be a valid
resource name.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Project.spec.parent.resourceRef
<sup><sup>[↩ Parent](#projectspecparent)</sup></sup>



Resource is a reference to the parent of the project. Must be a valid
resource.

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
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          Group is the group of the resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind is the kind of the resource.<br/>
          <br/>
            <i>Enum</i>: Organization<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of the resource.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Project.status
<sup><sup>[↩ Parent](#project)</sup></sup>



ProjectStatus defines the observed state of Project.

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
        <td><b><a href="#projectstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Represents the observations of a project's current state.
Known condition types are: "Ready"<br/>
          <br/>
            <i>Default</i>: [map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for control plane to reconcile reason:Unknown status:Unknown type:Ready]]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Project.status.conditions[index]
<sup><sup>[↩ Parent](#projectstatus)</sup></sup>



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
