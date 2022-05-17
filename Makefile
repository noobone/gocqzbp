##############################################################################
# Variables used for various build targets.
##############################################################################

# Enforce use of modules.
export GO111MODULE=on
export CGO_ENABLED=0

.PHONY: cqhttp
cqhttp:
	set -ex
	go build -ldflags "-s -w" -o cqhttp -trimpath