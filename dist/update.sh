#!/bin/sh

project=github.com/couchbaselabs/tuqtng
top=`go list -f '{{.Dir}}' $project`
version=`git describe`

cd $top

DIST=$top/dist

testpkg() {
    go test $project/...
    go vet $project/...
}

mkversion() {
    echo "{\"version\": \"$version\"}" > $DIST/version.json
}

build() {
    pkg=$project
    goflags="-v -ldflags '-X main.VERSION $version'"

    eval env GOARCH=386   GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq-engine.lin32 $pkg &
    eval env GOARCH=arm   GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq-engine.arm $pkg &
    eval env GOARCH=arm   GOARM=5 GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq-engine.arm5 $pkg &
    eval env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq-engine.lin64 $pkg &
    eval env GOARCH=amd64 GOOS=freebsd CGO_ENABLED=0 go build $goflags -o $DIST/cbq-engine.fbsd $pkg &&
    eval env GOARCH=386   GOOS=windows go build $goflags -o $DIST/cbq-engine.win32.exe $pkg &
    eval env GOARCH=amd64 GOOS=windows go build $goflags -o $DIST/cbq-engine.win64.exe $pkg &
    eval env GOARCH=amd64 GOOS=darwin go build $goflags -o $DIST/cbq-engine.mac $pkg &

    wait
}

buildclient() {
    pkg=$project/tuq_client
    goflags="-v -ldflags '-X main.VERSION $version'"

    eval env GOARCH=386   GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq.lin32 $pkg &
    eval env GOARCH=arm   GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq.arm $pkg &
    eval env GOARCH=arm   GOARM=5 GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq.arm5 $pkg &
    eval env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build $goflags -o $DIST/cbq.lin64 $pkg &
    eval env GOARCH=amd64 GOOS=freebsd CGO_ENABLED=0 go build $goflags -o $DIST/cbq.fbsd $pkg &
    eval env GOARCH=386   GOOS=windows go build $goflags -o $DIST/cbq.win32.exe $pkg &
    eval env GOARCH=amd64 GOOS=windows go build $goflags -o $DIST/cbq.win64.exe $pkg &
    eval env GOARCH=amd64 GOOS=darwin go build $goflags -o $DIST/cbq.mac $pkg &

    wait
}

builddistpackages() {

    mkdir -p $DIST/tutorial_tmp
    cd tutorial
    go build
    cd $top
    tutorial/tutorial -src tutorial/content/ -dst $DIST/tutorial_tmp/

    # mac build
    mkdir -p $DIST/stage
    cp $DIST/README $DIST/stage
    cp $DIST/license-ce-2013.txt $DIST/stage/LICENSE.txt
    cp $DIST/start_tutorial.sh $DIST/stage
    cp $DIST/cbq-engine.mac $DIST/stage/cbq-engine
    cp $DIST/cbq.mac $DIST/stage/cbq
    cp -r static/ $DIST/stage/static
    mkdir -p $DIST/stage/static/tutorial
    cp -r $DIST/tutorial_tmp/tutorial/content/ $DIST/stage/static/tutorial
    mkdir -p $DIST/stage/data/default/tutorial
    unzip tutorial/data/sampledb.zip -d $DIST/stage/data/default/
    cd $DIST/stage
    zip $DIST/couchbase-query_dev_preview2_x86_64_mac.zip -r .
    cd $top
    rm -rf $DIST/stage

    #linux 32
    mkdir -p $DIST/stage
    cp $DIST/README $DIST/stage
    cp $DIST/license-ce-2013.txt $DIST/stage/LICENSE.txt
    cp $DIST/start_tutorial.sh $DIST/stage
    cp $DIST/cbq-engine.lin32 $DIST/stage/cbq-engine
    cp $DIST/cbq.lin32 $DIST/stage/cbq
    cp -r static/ $DIST/stage/static
    mkdir -p $DIST/stage/static/tutorial
    cp -r $DIST/tutorial_tmp/tutorial/content/ $DIST/stage/static/tutorial
    mkdir -p $DIST/stage/data/default/tutorial
    unzip tutorial/data/sampledb.zip -d $DIST/stage/data/default/
    cd $DIST/stage
    tar zcvf $DIST/couchbase-query_dev_preview2_x86_linux.tar.gz .
    cd $top
    rm -rf $DIST/stage

    #linux 64
    mkdir -p $DIST/stage
    cp $DIST/README $DIST/stage
    cp $DIST/license-ce-2013.txt $DIST/stage/LICENSE.txt
    cp $DIST/start_tutorial.sh $DIST/stage
    cp $DIST/cbq-engine.lin64 $DIST/stage/cbq-engine
    cp $DIST/cbq.lin64 $DIST/stage/cbq
    cp -r static/ $DIST/stage/static
    mkdir -p $DIST/stage/static/tutorial
    cp -r $DIST/tutorial_tmp/tutorial/content/ $DIST/stage/static/tutorial
    mkdir -p $DIST/stage/data/default/tutorial
    unzip tutorial/data/sampledb.zip -d $DIST/stage/data/default/
    cd $DIST/stage
    tar zcvf $DIST/couchbase-query_dev_preview2_x86_64_linux.tar.gz .
    cd $top
    rm -rf $DIST/stage

    #win 32
    mkdir -p $DIST/stage
    cp $DIST/README $DIST/stage
    cp $DIST/license-ce-2013.txt $DIST/stage/LICENSE.txt
    cp $DIST/start_tutorial.bat $DIST/stage
    cp $DIST/cbq-engine.win32.exe $DIST/stage/cbq-engine.exe
    cp $DIST/cbq.win32.exe $DIST/stage/cbq.exe
    cp -r static/ $DIST/stage/static
    mkdir -p $DIST/stage/static/tutorial
    cp -r $DIST/tutorial_tmp/tutorial/content/ $DIST/stage/static/tutorial
    mkdir -p $DIST/stage/data/default/tutorial
    unzip tutorial/data/sampledb.zip -d $DIST/stage/data/default/
    cd $DIST/stage
    zip $DIST/couchbase-query_dev_preview2_x86_win.zip -r .
    cd $top
    rm -rf $DIST/stage

    #win 64
    mkdir -p $DIST/stage
    cp $DIST/README $DIST/stage
    cp $DIST/license-ce-2013.txt $DIST/stage/LICENSE.txt
    cp $DIST/start_tutorial.bat $DIST/stage
    cp $DIST/cbq-engine.win64.exe $DIST/stage/cbq-engine.exe
    cp $DIST/cbq.win64.exe $DIST/stage/cbq.exe
    cp -r static/ $DIST/stage/static
    mkdir -p $DIST/stage/static/tutorial
    cp -r $DIST/tutorial_tmp/tutorial/content/ $DIST/stage/static/tutorial
    mkdir -p $DIST/stage/data/default/tutorial
    unzip tutorial/data/sampledb.zip -d $DIST/stage/data/default/
    cd $DIST/stage
    zip $DIST/couchbase-query_dev_preview2_x86_64_win.zip -r .
    cd $top
    rm -rf $DIST/stage

    rm -rf $DIST/tutorial_tmp
    
}

compress() {
    rm -f $DIST/cbq-engine.*.gz $DIST/cbq.*.gz || true

    for i in $DIST/cbq-engine.* $DIST/cbq.*
    do
        gzip -9v $i &
    done

    wait
}

benchmark() {
    go test -test.bench . > $DIST/benchmark.txt
}

coverage() {
    for sub in ast misc plan test xpipeline
    do
        gocov test $project/$sub | gocov-html > $DIST/cov-$sub.html
    done
    cd $top/test
    gocov test -deps -exclude-goroot > $DIST/integ.json
    cat $DIST/integ.json | jq '{"Packages": [.Packages[] | if .Name > "github.com/couchbaselabs/tuqtng" and .Name < "github.com/couchbaselabs/tuqtnh" then . else empty end]}' > $DIST/integ2.json
    cat $DIST/integ2.json |gocov-html > $DIST/integ-cov.html
    cd $top
}

upload() {
    cbfsclient ${cbfsserver:-http://cbfs.hq.couchbase.com:8484/} upload \
        -ignore=$DIST/.cbfsclient.ignore -delete -v \
        $DIST/ tuqtng/

    cbfsclient ${cbfsserver:-http://cbfs.hq.couchbase.com:8484/} upload \
        -ignore=$DIST/.cbfsclient.ignore -delete -v \
        $DIST/redirect.html tuqtng
}

testpkg
mkversion
build
buildclient
builddistpackages
compress
benchmark
coverage
upload
