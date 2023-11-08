./node master topo_file.dat $1 &
./node backup1 topo_file.dat $1 &
./node backup2 topo_file.dat $1 &
./node backup3 topo_file.dat $1 &
./node witness1 topo_file.dat $1 &
while true; do sleep 86400; done