# Dynamic Resources Controller for kubernetes
**This application is still in early development**

A controller that is able to create (and update) kubenetes objects dynamically based on the status of other objects in the cluster. 
It works by generating the target resource and apply transformations on top of it.

**Note**

This controller is currently just a proof of concept and lacks a lot of its intended functionality and should be considered purely experimental

## Planned features
- Update resource on source resource change
- Advanced path-spec
- Advanced source resource spec (name-matchers, label- and field-matchers)
