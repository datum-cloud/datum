# API Reference

Packages:

- [iam.datumapis.com/v1alpha1](#iamdatumapiscomv1alpha1)

# iam.datumapis.com/v1alpha1

Resource Types:

- [GroupMembership](#groupmembership)




## GroupMembership
<sup><sup>[↩ Parent](#iamdatumapiscomv1alpha1 )</sup></sup>






GroupMembership is the Schema for the groupmemberships API

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
      <td>GroupMembership</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#groupmembershipspec">spec</a></b></td>
        <td>object</td>
        <td>
          GroupMembershipSpec defines the desired state of GroupMembership<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupmembershipstatus">status</a></b></td>
        <td>object</td>
        <td>
          GroupMembershipStatus defines the observed state of GroupMembership<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GroupMembership.spec
<sup><sup>[↩ Parent](#groupmembership)</sup></sup>



GroupMembershipSpec defines the desired state of GroupMembership

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
        <td><b><a href="#groupmembershipspecgroupref">groupRef</a></b></td>
        <td>object</td>
        <td>
          GroupRef is a reference to the Group.
Group is a namespaced resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#groupmembershipspecuserref">userRef</a></b></td>
        <td>object</td>
        <td>
          UserRef is a reference to the User that is a member of the Group.
User is a cluster-scoped resource.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GroupMembership.spec.groupRef
<sup><sup>[↩ Parent](#groupmembershipspec)</sup></sup>



GroupRef is a reference to the Group.
Group is a namespaced resource.

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
          Name is the name of the Group being referenced.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced Group.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GroupMembership.spec.userRef
<sup><sup>[↩ Parent](#groupmembershipspec)</sup></sup>



UserRef is a reference to the User that is a member of the Group.
User is a cluster-scoped resource.

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
          Name is the name of the User being referenced.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GroupMembership.status
<sup><sup>[↩ Parent](#groupmembership)</sup></sup>



GroupMembershipStatus defines the observed state of GroupMembership

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
        <td><b><a href="#groupmembershipstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions represent the latest available observations of an object's current state.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GroupMembership.status.conditions[index]
<sup><sup>[↩ Parent](#groupmembershipstatus)</sup></sup>



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
