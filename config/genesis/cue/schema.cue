package genesis

#User: {
	email:      !="" // must be specified and non-empty
	first_name: !="" // must be specified and non-empty
	last_name:  !="" // must be specified and non-empty
}

// Auth is the permissions and roles required for the Role Based Access Control (RBAC) setup of the app
#Auth: {
	permissions: [...#Permission]
	roles: [...#Role]
}

// Role is a job function or title which defines an authority level.
#Role: {
	// A human-readable code which represents the role.
	role_cd: !="" // must be specified and non-empty
	// A longer description of the role.
	role_description: !="" // must be specified and non-empty
	// A boolean denoting whether the role is active (true) or not (false).
	active: bool
	// A list of permissions that the role allows
	permissions: [...#Permission]
}

// Permission stores an approval of a mode of access to a resource.
#Permission: {
	// A human-readable string which represents a resource (e.g. an HTTP route or document, etc.).
	resource: !="" // must be specified and non-empty
	// A string representing the action taken on the resource (e.g. POST, GET, edit, etc.)
	operation: !="" // must be specified and non-empty
	// A description of what the permission is granting, e.g. "grants ability to edit a billing document".
	description: !="" // must be specified and non-empty
	// A boolean denoting whether the permission is active (true) or not (false).
	active: bool
}
