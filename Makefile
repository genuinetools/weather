# Setup name variables for the package/tool
NAME := weather
PKG := github.com/genuinetools/$(NAME)

CGO_ENABLED := 0

include basic.mk
