import os

def print_tree(directory, padding='|-- '):
    print(padding[:-3])
    padding += '   '
    for filename in sorted(os.listdir(directory)):
        path = os.path.join(directory, filename)
        if os.path.isdir(path):
            print(padding + filename + '/')
            print_tree(path, padding=padding + '|-- ')
        else:
            print(padding + filename)

if __name__ == "__main__":
    print_tree('.')