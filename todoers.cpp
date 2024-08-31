#include <iostream>
#include <sstream>
#include <string>
#include <map>
#include <vector>

// for string delimiter
std::vector<std::string> split(std::string s, std::string delimiter)
{
    size_t pos_start = 0, pos_end, delim_len = delimiter.length();
    std::string token;
    std::vector<std::string> res;

    while ((pos_end = s.find(delimiter, pos_start)) != std::string::npos)
    {
        token = s.substr(pos_start, pos_end - pos_start);
        pos_start = pos_end + delim_len;
        res.push_back(token);
    }

    res.push_back(s.substr(pos_start));
    return res;
}

std::vector<std::string> split(const std::string &s, char delim)
{
    std::vector<std::string> result;
    std::stringstream ss(s);
    std::string item;

    while (getline(ss, item, delim))
    {
        result.push_back(item);
    }

    return result;
}

std::string join(std::vector<std::string> str, char delim)
{
    std::string buf = "";
    for (std::string s : str)
    {
        buf += s;
        buf += delim;
    }

    return buf;
}

void add(std::string str) { std::cout << "add: " << str; }

void remove(std::string str) { std::cout << "remove: " << str; }

void move(std::string str) { std::cout << "move: " << str; }

using pFunc = void (*)(std::string);

// map of functions used
std::map<std::string, pFunc> funcMap{
    {"add", add},
    {"remove", remove},
    {"move", move}};

int main(int argc, char *argv[])
{

    if (argc > 1)
    {
        // get parameters into a vector
        std::vector<std::string> args(argv, argv + argc);

        // remove the first element - the executable
        args.erase(args.begin());

        // get the command from the first index
        std::string cmd = args[0];

        // remove the command element
        args.erase(args.begin());

        // join the rest as the passed values
        std::string value = join(args, ' ');

        // process the input based on the command
        funcMap[cmd](value);
    }
    else
    {
        std::cout << "usage:\ttodoers <cmd> <...values>" << std::endl
                  << std::endl
                  << "commands: " << std::endl;

        for (auto const &[key, val] : funcMap)
        {
            std::cout << "\t" << key << std::endl;
        }
    }

    return 0;
}