CMAKE_MINIMUM_REQUIRED (VERSION 2.8)
PROJECT (tuqtng)

# This project is intended to be built as part of a Couchbase Server build,
# but for nominal support building stand-alone, the following variables
# need to be specified:
#
# CMAKE_MODULE_PATH - should contain a directory containing
#      FindCouchbaseGo.cmake, eg. tlm/cmake/Modules
#
# CMAKE_INSTALL_PREFIX - should point to a directory where the bin/ output
#      directory will be placed
#
# GODEPSPATH - should point to a Go workspace directory containing all
#      transitive Go dependencies

INCLUDE (FindCouchbaseGo)

GoInstall (TARGET cbq PACKAGE github.com/couchbaselabs/tuqtng/cbq
  GOPATH "${PROJECT_SOURCE_DIR}/../../../.." "${GODEPSPATH}"
  INSTALL_PATH bin)

GoInstall (TARGET cbq-engine PACKAGE github.com/couchbaselabs/tuqtng
  GOPATH "${PROJECT_SOURCE_DIR}/../../../.." "${GODEPSPATH}"
  INSTALL_PATH bin OUTPUT cbq-engine)
