#include <iostream>
#include <unordered_map>
#include <vector>

#define Atom char
#define PathVector std::vector<char>
#define ChildrenMap std::unordered_map<char, Node*>
#define MininodeMap std::unordered_map<std::string, Atom*>

struct Node {
  ChildrenMap children;
  MininodeMap mininodes;
  PathVector path;

  Node () { };
  Node (char c) {
    Atom a = c;
    if (children.find(c) != children.end()) {
      throw std::runtime_error("Duplication of chars");
    }
  };
};

class TreeDoc {
 private:
  Node* root = nullptr;

 public:
  TreeDoc() {
    root = new Node();  // Root node representing an empty character
  }

  void insert(char c, PathVector path) {
    Node* curr = root;
    for (char direction : path) {
      if (curr->children.find(direction) == curr->children.end()) {
        curr->children[direction] = new Node();
      }
      curr = curr->children[direction];
    }
    
    Atom* newAtom = new Atom(c);
    std::string posId = generatePosId(curr);
    curr->mininodes[posId] = newAtom;
    curr->path = path;
  }

  // Function to traverse and print the document
  void printDocument() {
    traverse(root);
  }

 private:
  // Helper function for traversing the tree
  void traverse(Node* node) {
    if (!node) {
      return;
    }

    std::cout << "PATH: ";
    for (const auto& a : node->path) {
      std::cout << a << " ";
    }
    std::cout << std::endl;
    
    std::cout << "ELEMENTS: ";
    for (const auto& a : node->mininodes) {
      std::cout << a.second << " ";
    }
    std::cout << std::endl << std::endl;

    for (const auto& child : node->children) {
      traverse(child.second);
    }
  }

  std::string generatePosId(Node* node) {
    size_t mininode_ct = node->mininodes.size();
    return std::to_string(mininode_ct);
  }
};

int main() {
  TreeDoc document;

  // Example insertions
  document.insert('a', {'L', 'R'});
  document.insert('b', {'L', 'L'});
  document.insert('c', {'R', 'L'});
  document.insert('d', {'R', 'R'});
  document.insert('e', {'R', 'R'});

  document.printDocument();  // Output the current document content

  return 0;
}
