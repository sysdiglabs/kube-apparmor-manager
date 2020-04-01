package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *ApparmorProfile) DeepCopyInto(out *ApparmorProfile) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = ApparmorProfileSpec{
		Rules:    in.Spec.Rules,
		Enforced: in.Spec.Enforced,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ApparmorProfile) DeepCopyObject() runtime.Object {
	out := ApparmorProfile{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ApparmorProfileList) DeepCopyObject() runtime.Object {
	out := ApparmorProfileList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]ApparmorProfile, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
