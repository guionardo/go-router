package router

func (r *Router) assertInfo() *Router {
	if r.info == nil {
		r.info = &RouterInfo{}
	}
	return r
}

func Title(title string) func(*Router) {
	return func(r *Router) {
		r.assertInfo().info.Title = title
	}
}

func Version(version string) func(*Router) {
	return func(r *Router) {
		r.assertInfo().info.Version = version
	}
}

func Description(description string) func(*Router) {
	return func(r *Router) {
		r.assertInfo().info.Description = description
	}
}
