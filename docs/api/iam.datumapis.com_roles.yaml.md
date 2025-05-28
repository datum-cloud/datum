# API Reference

Packages:

- [iam.datumapis.com/v1alpha1](#iamdatumapiscomv1alpha1)

# iam.datumapis.com/v1alpha1

Resource Types:

- [Role](#role)




## Role
<sup><sup>[↩ Parent](#iamdatumapiscomv1alpha1 )</sup></sup>






Role is the Schema for the roles API

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
      <td>Role</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#rolespec">spec</a></b></td>
        <td>object</td>
        <td>
          RoleSpec defines the desired state of Role<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#rolestatus">status</a></b></td>
        <td>object</td>
        <td>
          RoleStatus defines the observed state of Role<br/>
          <br/>
            <i>Default</i>: map[conditions:[map[lastTransitionTime:1970-01-01T00:00:00Z message:Waiting for control plane to reconcile reason:Unknown status:Unknown type:Ready]]]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Role.spec
<sup><sup>[↩ Parent](#role)</sup></sup>



RoleSpec defines the desired state of Role

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
        <td><b>includedPermissions</b></td>
        <td>[]string</td>
        <td>
          The names of the permissions this role grants when bound in an IAM policy.
All permissions must be in the format: `{service}.{resource}.{action}`
(e.g. compute.workloads.create).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>launchStage</b></td>
        <td>string</td>
        <td>
          Defines the launch stage of the IAM Role. Must be one of: Early Access,
Alpha, Beta, Stable, Deprecated.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#rolespecinheritedrolesindex">inheritedRoles</a></b></td>
        <td>[]object</td>
        <td>
          The list of roles from which this role inherits permissions.
Each entry must be a valid role resource name.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Role.spec.inheritedRoles[index]
<sup><sup>[↩ Parent](#rolespec)</sup></sup>



ScopedRoleReference defines a reference to another Role, scoped by namespace.
This is used for purposes like role inheritance where a simple name and namespace
is sufficient to identify the target role.

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
          Name of the referenced Role.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced Role.
If not specified, it defaults to the namespace of the resource containing this reference.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Role.status
<sup><sup>[↩ Parent](#role)</sup></sup>



RoleStatus defines the observed state of Role

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
        <td><b><a href="#rolestatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions provide conditions that represent the current status of the Role.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          ObservedGeneration is the most recent generation observed by the controller.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parent</b></td>
        <td>string</td>
        <td>
          The resource name of the parent the role was created under.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Role.status.conditions[index]
<sup><sup>[↩ Parent](#rolestatus)</sup></sup>



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
