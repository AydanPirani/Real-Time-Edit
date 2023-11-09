#ifndef ATOMBUFFER_H
#define ATOMBUFFER_H

#include <map>
#include <vector>

#define PosID int

struct Atom {
    char content;
};

// Declare the AtomBuffer class
class AtomBuffer {
private:
    std::map<PosID, Atom> atomMap;

public:
    // Constructor
    AtomBuffer();

    // Member functions
    int Insert(PosID, Atom);
    int Delete(PosID);

    // Helper/util functions
    void PrintContents();
    std::vector<PosID> GetKeys();

};

#endif