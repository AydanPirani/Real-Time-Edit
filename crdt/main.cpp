#include <iostream>
#include <unordered_map>
#include <vector>

#define Atom char
#define PathVector std::vector<char>
#define ChildrenMap std::unordered_map<char, Node*>
#define MininodeMap std::unordered_map<std::string, Atom>

struct Node {
  ChildrenMap children;
  MininodeMap mininodes;
  PathVector path;

  Node () { };
  Node (Atom a) {
    if (children.find(a) != children.end()) {
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

  // Lazily insert the atom into the tree at the given path, while adding nodes to accomodate
  void insert(Atom a, PathVector path) {
    Node* curr = root;
    for (char direction : path) {
      
      // No child here -> create an empty path
      if (curr->children.find(direction) == curr->children.end()) {
        curr->children[direction] = new Node();
      }
      curr = curr->children[direction];
    }
    
    std::string posId = generateDisambiguator(curr);
    curr->mininodes[posId] = a;
    curr->path = path;
  }

  // Taverse and print the document
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

  // Function to generate the internal disambiguator - currently at just the Id
  std::string generateDisambiguator(Node* node) {
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
