./node master topo_file.dat $1 master_pipe &
./node backup1 topo_file.dat $1 backup1_pipe &
./node backup2 topo_file.dat $1 backup2_pipe &
./node backup3 topo_file.dat $1 backup3_pipe &
./node witness1 topo_file.dat $1 &
while true; do sleep 86400; done