#include <iostream>
#include <sstream>
#include <vector>

// for string delimiter
std::vector<std::string> split(std::string s, std::string delimiter) {
    size_t pos_start = 0, pos_end, delim_len = delimiter.length();
    std::string token;
    std::vector<std::string> res;

    while ((pos_end = s.find(delimiter, pos_start)) != std::string::npos) {
        token = s.substr (pos_start, pos_end - pos_start);
        pos_start = pos_end + delim_len;
        res.push_back (token);
    }

    res.push_back (s.substr (pos_start));
    return res;
}

std::vector<std::string> split (const std::string &s, char delim) {
    std::vector<std::string> result;
    std::stringstream ss (s);
    std::string item;

    while (getline (ss, item, delim)) {
        result.push_back (item);
    }

    return result;

}

std::string join(std::vector<std::string> str, char delim) {
	std::string buf = "";
	for(std::string s : str) {
		buf += s;
		buf += delim;
	}

	return buf;
}

void add() {

}

void remove() {

}

void move() {

}

int main(int argc, char *argv[]) {
	//get parameters into a vector
	std::vector<std::string> args(argv, argv + argc);
	//remove the first element - the executable
	args.erase(args.begin());
	
	std::cout << "Split: " << std::endl;
	for (std::string arg : args) std::cout << arg << std::endl;

	std::cout << args[0] << std::endl;

	switch(args[0]) {
		case "add":
			add();
			break;
	}

	std::cout << std::endl << "Joined: " << std::endl;
	std::string joinedArgs = join(args, ' ');

	std::cout << joinedArgs;

    return 0;
}

