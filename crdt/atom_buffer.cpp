#include "atom_buffer.h"
#include "iostream"

// Implementation of AtomBuffer constructor
AtomBuffer::AtomBuffer() {
    // Do nothing except initialize the map
    AtomBuffer::atomMap = std::map<PosID, Atom>();
}

// Print everything in the map
void AtomBuffer::PrintContents() {
    for (const auto& entry : atomMap) {
        std::cout << "PosID: " << entry.first << ", Content: " << entry.second.content << std::endl;
    }
}

// Insert element into the map
// Returns 0 if successful, else 1 if pos exists
int AtomBuffer::Insert(PosID pos, Atom atom) {
    // Insert the atom into the map, assuming PosID is unique

    if (atomMap.find(pos) != atomMap.end()) {
        return 1;
    }

    atomMap[pos] = atom;
    return 0;
}

// Remove element from the map
// Returns 0 if deleted, else 1
int AtomBuffer::Delete(PosID pos) {
    int deleted = atomMap.erase(pos);
    return (deleted == 1) ? 0 : 1;
}

std::vector<PosID> AtomBuffer::GetKeys() {
    std::vector<PosID> keys;

    // Iterate through the map and push keys into the vector
    for (const auto& pair : atomMap) {
        keys.push_back(pair.first);
    }
    return keys;
}