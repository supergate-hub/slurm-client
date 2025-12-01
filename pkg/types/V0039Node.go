// SPDX-FileCopyrightText: Copyright (C) SchedMD LLC.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"k8s.io/utils/ptr"
	"k8s.io/utils/set"

	api "github.com/supergate-hub/slurm-client/api/v0039"
	"github.com/supergate-hub/slurm-client/pkg/object"
	"github.com/supergate-hub/slurm-client/pkg/utils"
)

const (
	ObjectTypeV0039Node = "V0039Node"
)

type V0039Node struct {
	api.V0039Node
}

// GetKey implements Object.
func (o *V0039Node) GetKey() object.ObjectKey {
	return object.ObjectKey(ptr.Deref(o.Name, ""))
}

// GetType implements Object.
func (o *V0039Node) GetType() object.ObjectType {
	return ObjectTypeV0039Node
}

// DeepCopyObject implements Object.
func (o *V0039Node) DeepCopyObject() object.Object {
	return o.DeepCopy()
}

func (o *V0039Node) DeepCopy() *V0039Node {
	out := new(V0039Node)
	utils.RemarshalOrDie(o, out)
	return out
}

func (o *V0039Node) GetStateAsSet() set.Set[api.V0039NodeState] {
	out := make(set.Set[api.V0039NodeState])
	states := ptr.Deref(o.State, []api.V0039NodeState{})
	for _, s := range states {
		out.Insert(s)
	}
	return out
}

type V0039NodeList struct {
	Items []V0039Node
}

// GetType implements ObjectList.
func (o *V0039NodeList) GetType() object.ObjectType {
	return ObjectTypeV0039Node
}

// GetItems implements ObjectList.
func (o *V0039NodeList) GetItems() []object.Object {
	list := make([]object.Object, len(o.Items))
	for i, item := range o.Items {
		list[i] = item.DeepCopyObject()
	}
	return list
}

// AppendItem implements ObjectList.
func (o *V0039NodeList) AppendItem(object object.Object) {
	out, ok := object.(*V0039Node)
	if ok {
		utils.RemarshalOrDie(object, out)
		o.Items = append(o.Items, *out)
	}
}

// DeepCopyObjectList implements ObjectList.
func (o *V0039NodeList) DeepCopyObjectList() object.ObjectList {
	out := new(V0039NodeList)
	out.Items = make([]V0039Node, len(o.Items))
	for i, item := range o.Items {
		out.Items[i] = *item.DeepCopy()
	}
	return out
}

