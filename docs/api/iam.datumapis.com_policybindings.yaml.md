# API Reference

Packages:

- [iam.datumapis.com/v1alpha1](#iamdatumapiscomv1alpha1)

# iam.datumapis.com/v1alpha1

Resource Types:

- [PolicyBinding](#policybinding)




## PolicyBinding
<sup><sup>[↩ Parent](#iamdatumapiscomv1alpha1 )</sup></sup>






PolicyBinding is the Schema for the policybindings API

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
      <td>PolicyBinding</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#policybindingspec">spec</a></b></td>
        <td>object</td>
        <td>
          PolicyBindingSpec defines the desired state of PolicyBinding<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#policybindingstatus">status</a></b></td>
        <td>object</td>
        <td>
          PolicyBindingStatus defines the observed state of PolicyBinding<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PolicyBinding.spec
<sup><sup>[↩ Parent](#policybinding)</sup></sup>



PolicyBindingSpec defines the desired state of PolicyBinding

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
        <td><b><a href="#policybindingspecroleref">roleRef</a></b></td>
        <td>object</td>
        <td>
          RoleRef is a reference to the Role that is being bound.
This can be a reference to a Role custom resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#policybindingspecsubjectsindex">subjects</a></b></td>
        <td>[]object</td>
        <td>
          Subjects holds references to the objects the role applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#policybindingspectargetref">targetRef</a></b></td>
        <td>object</td>
        <td>
          TargetRef is a reference to the resource to which this policy binding applies.
This allows the binding to be about a resource in any namespace or a cluster-scoped resource.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### PolicyBinding.spec.roleRef
<sup><sup>[↩ Parent](#policybindingspec)</sup></sup>



RoleRef is a reference to the Role that is being bound.
This can be a reference to a Role custom resource.

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
          APIGroup is the group for the resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced Role. If empty, it is assumed to be in the PolicyBinding's namespace.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PolicyBinding.spec.subjects[index]
<sup><sup>[↩ Parent](#policybindingspec)</sup></sup>



Subject contains a reference to the object or user identities a role binding applies to.
This can be a User or Group.

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
          APIGroup holds the API group of the referenced subject.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of object being referenced. Values defined in Kind constants.<br/>
          <br/>
            <i>Enum</i>: User<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the object being referenced.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>uid</b></td>
        <td>string</td>
        <td>
          UID of the referenced object.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced object. If DNE, then for an SA it refers to the PolicyBinding resource's namespace.
For a User or Group, it is ignored.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PolicyBinding.spec.targetRef
<sup><sup>[↩ Parent](#policybindingspec)</sup></sup>



TargetRef is a reference to the resource to which this policy binding applies.
This allows the binding to be about a resource in any namespace or a cluster-scoped resource.

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
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>uid</b></td>
        <td>string</td>
        <td>
          UID is the unique identifier of the resource being referenced.<br/>
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
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace is the namespace of resource being referenced.
Required for namespace-scoped resources. Omitted for cluster-scoped resources.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PolicyBinding.status
<sup><sup>[↩ Parent](#policybinding)</sup></sup>



PolicyBindingStatus defines the observed state of PolicyBinding

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
        <td><b><a href="#policybindingstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions provide conditions that represent the current status of the PolicyBinding.<br/>
          <br/>
            <i>Default</i>: [map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for control plane to reconcile reason:Unknown status:Unknown type:Ready]]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the most recent generation observed for this PolicyBinding by the controller.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PolicyBinding.status.conditions[index]
<sup><sup>[↩ Parent](#policybindingstatus)</sup></sup>



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
