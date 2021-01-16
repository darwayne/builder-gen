// Code generated by builder-gen. DO NOT EDIT.
package generator

func (o DirOpts) HasRecursiveExclusions() bool {
	return o.RecursiveExclusions != nil
}

func (o DirOpts) ToOptFuncs() []DirOptsFunc {
	builder := NewDirOptsBuilder()

	if o.HasRecursiveExclusions() {
		builder.RecursiveExclusions(o.RecursiveExclusions...)
	}

	builder.Recursive(o.Recursive)

	return builder.Build()
}

type DirOptsFunc func(*DirOpts)
type DirOptsBuilder struct {
	opts []DirOptsFunc
}

func NewDirOptsBuilder(opts ...DirOptsFunc) *DirOptsBuilder {
	builder := &DirOptsBuilder{opts: opts}
	return builder
}

func (l *DirOptsBuilder) Recursive(recursive bool) *DirOptsBuilder {
	return l.add(func(opts *DirOpts) {
		opts.Recursive = recursive
	})
}
func (l *DirOptsBuilder) RecursiveExclusions(recursiveExclusions ...string) *DirOptsBuilder {
	return l.add(func(opts *DirOpts) {
		opts.RecursiveExclusions = recursiveExclusions
	})
}

func (l *DirOptsBuilder) add(fn DirOptsFunc) *DirOptsBuilder {
	l.opts = append(l.opts, fn)
	return l
}

func (l *DirOptsBuilder) Build() []DirOptsFunc {
	return l.opts
}

func ToDirOpts(opts ...DirOptsFunc) DirOpts {
	var info DirOpts
	ToDirOptsWithDefault(&info, opts...)

	return info
}

func ToDirOptsWithDefault(info *DirOpts, opts ...DirOptsFunc) {
	for _, o := range opts {
		o(info)
	}
}