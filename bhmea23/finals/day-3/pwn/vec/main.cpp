#include <stdexcept>
#include <iostream>
#include <vector>

int main() {
  std::vector<size_t> vec;

  std::setbuf(stdin, nullptr);
  std::setbuf(stdout, nullptr);
  std::setbuf(stderr, nullptr);
  std::cout << "1. set" << std::endl
            << "2. get" << std::endl
            << "3. copy" << std::endl
            << "4. clear" << std::endl;

  while (std::cin.good()) {
    size_t choice;
    std::cout << ">> ";
    std::cin >> choice;

    switch (choice) {
      case 1: {
        /* set */
        size_t index, value;
        std::cout << "index: ";
        std::cin >> index;
        if (index > vec.size())
          throw std::out_of_range("vector index out of range");

        std::cout << "value: ";
        std::cin >> value;
        if (index < vec.size())
          vec[index] = value;
        else if (index == vec.size())
          vec.emplace_back(value);
        break;
      }

      case 2: {
        /* get */
        size_t index, value;
        std::cout << "index: ";
        std::cin >> index;
        if (index >= vec.size())
          throw std::out_of_range("vector index out of range");

        std::cout << "vec[" << index << "] = " << vec[index] << std::endl;
        break;
      }

      case 3: {
        /* copy */
        size_t src, dest, count;
        std::cout << "from: ";
        std::cin >> src;
        std::cout << "to: ";
        std::cin >> dest;
        std::cout << "count: ";
        std::cin >> count;

        if (src > vec.size() || dest > vec.size())
          throw std::out_of_range("vector index out of range");
        if (src + count > vec.size() || dest + count > vec.size())
          throw std::out_of_range("count too big");
        std::copy(vec.begin() + src,
                  vec.begin() + src + count,
                  vec.begin() + dest);
        break;
      }

      case 4:
        /* clear */
        vec.clear();
        vec.shrink_to_fit();
        break;

      default:
        return 0;
    }
  }

  return 1;
}
