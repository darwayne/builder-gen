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
	builder.Trace(o.Trace)

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

func (l *DirOptsBuilder) Recursive(recursiveParam bool) *DirOptsBuilder {
	return l.add(func(opts *DirOpts) {
		opts.Recursive = recursiveParam
	})
}
func (l *DirOptsBuilder) RecursiveExclusions(recursiveExclusionsParam ...string) *DirOptsBuilder {
	return l.add(func(opts *DirOpts) {
		opts.RecursiveExclusions = recursiveExclusionsParam
	})
}
func (l *DirOptsBuilder) Trace(traceParam bool) *DirOptsBuilder {
	return l.add(func(opts *DirOpts) {
		opts.Trace = traceParam
	})
}

func (l *DirOptsBuilder) add(fn DirOptsFunc) *DirOptsBuilder {
	l.opts = append(l.opts, fn)
	return l
}

func (l *DirOptsBuilder) Build() []DirOptsFunc {
	return l.opts
}

func (l *DirOptsBuilder) ToDirOpts() DirOpts {
	return ToDirOpts(l.opts...)
}

func (l *DirOptsBuilder) ToDirOptsWithDefault(info *DirOpts) {
	ToDirOptsWithDefault(info, l.opts...)
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
