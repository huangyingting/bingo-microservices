API_POST_PROCESS=cd ../../api/$(APP_NAME) && ls v1/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
include ../../Makefile.inc