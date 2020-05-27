#!/usr/bin/env python3

import os
import argparse

def test():
    print ("-- vet broker")
    os.system("go vet ./...")
    print ("-- test broker")
    os.system("go test -v ./...")

def build():
    print ("-- clean broker")
    os.system("go clean -r -cache -testcache -modcache")
    print ("-- fmt broker")
    os.system("go fmt ./...")
    print("-- build broker")
    print("go build " + ldflags() + " -v ./...")
    os.system("go build " + ldflags() + " -v ./...")

def run():
    print ("-- run broker")
    os.system("go run  " + ldflags() + " brokerApp.go")

def dispatcher(cmd):
   dispatcher={
       'build': build,
       "run": run,
       "test": test,
   }
   return dispatcher.get(cmd)
 
def ldflags():
    dirtyStream = os.popen('git diff --quiet || echo dirty')
    dirty = dirtyStream.read()
    if dirty == "":
      commitStream = os.popen('git rev-parse HEAD')
      commit = commitStream.read()
    else:
      commit = dirty

    # TDOD: version (consider goreleaser) 
    return F"-ldflags=\"-X 'main.Version=dev' -X 'main.Commit={commit}'\""

def main():
  parser = argparse.ArgumentParser(description="Make tool for cloud foundry api broker", epilog="(c) 2020 by KLÄFF-Soft)")
  parser.add_help=True
  parser.add_argument("command", nargs='?', choices=['build', 'run', 'test'], help="commands to execute")
 
  args = parser.parse_args()
  func = dispatcher(args.command)

  if func != None:
    func()
  else: 
    parser.print_usage()

if __name__ == '__main__':
    main()