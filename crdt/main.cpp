#include <iostream>
#include <unordered_map>
#include <vector>

#define Atom char
#define Path std::vector<int>
#define Disambiguator int

#define ChildrenMap std::unordered_map<Disambiguator, std::unordered_map<int, Node *>>
#define ValueMap std::unordered_map<Disambiguator, Atom>

struct PosId {
    Path path;
    std::vector<Disambiguator> disambiguator;
};

struct Node {
    ChildrenMap children;
    ValueMap values;

    PosId posId;
    Node *parent;

    Node(){};
    Node(Atom a) {
        if (children.find(a) != children.end()) {
            throw std::runtime_error("Duplication of chars");
        }
    };
};

class TreeDoc {
   private:
    Node *root = nullptr;

   public:
    TreeDoc() {
        root = new Node();  // Root node representing an empty character
    }

    // Lazily insert the atom into the tree at the given path, while adding nodes
    // to accomodate
    void insert(Atom a, PosId posId) {
        Node *curr = getNodeAtPosId(posId);
        Disambiguator disambiguator = posId.disambiguator.back();

        curr->posId = posId;
        curr->values[disambiguator] = a;
    }

    // Taverse and print the document
    void printDocument() { traverse(root); }

   private:
    // Helper function for traversing the tree
    void traverse(Node *node) {
        if (!node) {
            return;
        }

        if (node->values.size() != 0) {
            std::cout << "PATH: ";
            
            // for (const auto &a : node->posId.path) {
            //     std::cout << a << " ";
            // }

            for (size_t i = 0; i < node->posId.path.size(); i++) {
              std::cout << "(" << node->posId.path[i] << "," << node->posId.disambiguator[i] << ") ";
            }

            std::cout << std::endl;

            std::cout << "ELEMENTS: ";
            for (const auto &a : node->values) {
                std::cout << a.second << " ";
            }
            std::cout << std::endl << std::endl;
        }

        for (const auto &mininode : node->children) {
            for (const auto &child : mininode.second) {
                traverse(child.second);
            }
        }
    }

    // Function to generate the internal disambiguator - currently at just the Id
    Disambiguator generateDisambiguator(Node *node) {
        int mininode_ct = node->values.size();
        return mininode_ct;
    }

    Node *getNodeAtPosId(PosId &posId) {
        Node *curr = root;

        Path runningPath;
        std::vector<Disambiguator> runningDisambiguator;

        for (size_t i = 0; i < posId.path.size(); i++) {
            int direction = posId.path[i];
            Disambiguator disambiguator = posId.disambiguator[i];

            runningPath.push_back(direction);
            runningDisambiguator.push_back(disambiguator);

            // No child here -> create an empty Node here
            if (curr->children[disambiguator].find(direction) == curr->children[disambiguator].end()) {
                Node *newNode = new Node();
                curr->children[disambiguator][direction] = newNode;
                newNode->parent = curr;
                newNode->posId = {.path = runningPath, .disambiguator = runningDisambiguator};
            }

            curr = curr->children[disambiguator][direction];
        }
        return curr;
    }

    PosId generatePosId(PosId p, PosId f, Disambiguator d) {
        // SKIPPING CHECK THAT THERE'S NO ATOM BETWEEN P AND F
        bool foundM = false;

        Node *pNode = getNodeAtPosId(p), *fNode = getNodeAtPosId(f);
        for (const auto& node : pNode->values) {
          if (node.first > p.disambiguator.back() && fNode->parent == pNode) {
            foundM = true;
            break;
          }
        }
        
        PosId newPosId;        
        if (isAncestorOf(p, f)) {
            newPosId.path = f.path;
            newPosId.path.push_back(0);
            newPosId.disambiguator = f.disambiguator;
            newPosId.disambiguator.push_back(d);
        } else if (isAncestorOf(f, p)) {
            newPosId.path = p.path;
            newPosId.path.push_back(1);
            newPosId.disambiguator = p.disambiguator;
            newPosId.disambiguator.push_back(d);
        } else if (areMiniSiblings(p, f) || foundM) {
            newPosId.path = p.path;
            newPosId.path.push_back(1);
            newPosId.disambiguator = p.disambiguator;
            newPosId.disambiguator.push_back(d);
        } else {
            newPosId.path = p.path;
            newPosId.path.push_back(1);
            newPosId.disambiguator = p.disambiguator;
            newPosId.disambiguator.push_back(d);
        }
        return newPosId;
    }

    bool isAncestorOf(PosId &expectedParent, PosId &expectedChild) {
        Path &parentPath = expectedParent.path, &childPath = expectedChild.path;
        std::vector<Disambiguator> &parentDisambiguator = expectedParent.disambiguator;
        std::vector<Disambiguator> &childDisambiguator = expectedChild.disambiguator;

        if (parentPath.size() > childPath.size()) {
            return false;
        }

        for (size_t i = 0; i < parentPath.size(); i++) {
            // Check that all actual path indices are the same, up until the parent node isreached
            if (parentPath[i] != childPath[i]) {
                return false;
            }

            // Check that all disambiguators are the same, up until the parent node is reached
            if (parentDisambiguator[i] != childDisambiguator[i]) {
                return false;
            }
        }

        return true;
    }

    // Mini-siblings IFF both paths are the same, and disambiguators are all the
    // same except the last one
    bool areMiniSiblings(PosId &u, PosId &v) {
        Path &uPath = u.path, &vPath = v.path;
        std::vector<Disambiguator> &uDisambiguator = u.disambiguator, &vDisambiguator = v.disambiguator;

        if (uDisambiguator.size() != vDisambiguator.size()) {
            return false;
        }

        // Check that all disambiguators are the same, up until very last node
        for (size_t i = 0; i < uDisambiguator.size() - 1; i++) {
            if (uDisambiguator[i] != vDisambiguator[i]) {
                return false;
            }
        }

        return u.path == v.path;
    }
};

int main() {
    TreeDoc document;

    document.insert('a', {{0, 1}, {0, 0}});
    document.insert('b', {{0, 1}, {0, 1}});
    document.insert('B', {{0, 1}, {0, 1}});
    document.insert('c', {{0}, {0}});
    document.insert('d', {{0, 1}, {0, 0}});
    // document.insert('e', {1, 0});

    document.printDocument();  // Output the current document content

    return 0;
}
