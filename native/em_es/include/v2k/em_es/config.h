//  To parse this JSON data, first install
//
//      json.hpp  https://github.com/nlohmann/json
//
//  Then include this file, and then do
//
//     Config data = nlohmann::json::parse(jsonString);

#pragma once

#include <nlohmann/json.hpp>

#include <optional>
#include <stdexcept>
#include <regex>

namespace v2k {
namespace em_es {
    using nlohmann::json;

    #ifndef NLOHMANN_UNTYPED_v2k_em_es_HELPER
    #define NLOHMANN_UNTYPED_v2k_em_es_HELPER
    inline json get_untyped(const json & j, const char * property) {
        if (j.find(property) != j.end()) {
            return j.at(property).get<json>();
        }
        return json();
    }

    inline json get_untyped(const json & j, std::string property) {
        return get_untyped(j, property.data());
    }
    #endif

    class X {
        public:
        X() = default;
        virtual ~X() = default;

        private:
        double weight;

        public:
        const double & get_weight() const { return weight; }
        double & get_mutable_weight() { return weight; }
        void set_weight(const double & value) { this->weight = value; }
    };

    class Axes {
        public:
        Axes() = default;
        virtual ~Axes() = default;

        private:
        X x;
        X z;

        public:
        const X & get_x() const { return x; }
        X & get_mutable_x() { return x; }
        void set_x(const X & value) { this->x = value; }

        const X & get_z() const { return z; }
        X & get_mutable_z() { return z; }
        void set_z(const X & value) { this->z = value; }
    };

    class V2KSchema {
        public:
        V2KSchema() = default;
        virtual ~V2KSchema() = default;

        private:
        Axes axes;

        public:
        const Axes & get_axes() const { return axes; }
        Axes & get_mutable_axes() { return axes; }
        void set_axes(const Axes & value) { this->axes = value; }
    };

    class Solver {
        public:
        Solver() = default;
        virtual ~Solver() = default;

        private:
        double absolute_cost;
        double max_evaluations;
        double max_iterations;
        double relative_cost;

        public:
        const double & get_absolute_cost() const { return absolute_cost; }
        double & get_mutable_absolute_cost() { return absolute_cost; }
        void set_absolute_cost(const double & value) { this->absolute_cost = value; }

        const double & get_max_evaluations() const { return max_evaluations; }
        double & get_mutable_max_evaluations() { return max_evaluations; }
        void set_max_evaluations(const double & value) { this->max_evaluations = value; }

        const double & get_max_iterations() const { return max_iterations; }
        double & get_mutable_max_iterations() { return max_iterations; }
        void set_max_iterations(const double & value) { this->max_iterations = value; }

        const double & get_relative_cost() const { return relative_cost; }
        double & get_mutable_relative_cost() { return relative_cost; }
        void set_relative_cost(const double & value) { this->relative_cost = value; }
    };

    class Config {
        public:
        Config() = default;
        virtual ~Config() = default;

        private:
        std::map<std::string, V2KSchema> cams;
        Solver solver;
        double window;

        public:
        const std::map<std::string, V2KSchema> & get_cams() const { return cams; }
        std::map<std::string, V2KSchema> & get_mutable_cams() { return cams; }
        void set_cams(const std::map<std::string, V2KSchema> & value) { this->cams = value; }

        const Solver & get_solver() const { return solver; }
        Solver & get_mutable_solver() { return solver; }
        void set_solver(const Solver & value) { this->solver = value; }

        const double & get_window() const { return window; }
        double & get_mutable_window() { return window; }
        void set_window(const double & value) { this->window = value; }
    };
}
}

namespace v2k {
namespace em_es {
    void from_json(const json & j, X & x);
    void to_json(json & j, const X & x);

    void from_json(const json & j, Axes & x);
    void to_json(json & j, const Axes & x);

    void from_json(const json & j, V2KSchema & x);
    void to_json(json & j, const V2KSchema & x);

    void from_json(const json & j, Solver & x);
    void to_json(json & j, const Solver & x);

    void from_json(const json & j, Config & x);
    void to_json(json & j, const Config & x);

    inline void from_json(const json & j, X& x) {
        x.set_weight(j.at("weight").get<double>());
    }

    inline void to_json(json & j, const X & x) {
        j = json::object();
        j["weight"] = x.get_weight();
    }

    inline void from_json(const json & j, Axes& x) {
        x.set_x(j.at("x").get<X>());
        x.set_z(j.at("z").get<X>());
    }

    inline void to_json(json & j, const Axes & x) {
        j = json::object();
        j["x"] = x.get_x();
        j["z"] = x.get_z();
    }

    inline void from_json(const json & j, V2KSchema& x) {
        x.set_axes(j.at("axes").get<Axes>());
    }

    inline void to_json(json & j, const V2KSchema & x) {
        j = json::object();
        j["axes"] = x.get_axes();
    }

    inline void from_json(const json & j, Solver& x) {
        x.set_absolute_cost(j.at("absoluteCost").get<double>());
        x.set_max_evaluations(j.at("maxEvaluations").get<double>());
        x.set_max_iterations(j.at("maxIterations").get<double>());
        x.set_relative_cost(j.at("relativeCost").get<double>());
    }

    inline void to_json(json & j, const Solver & x) {
        j = json::object();
        j["absoluteCost"] = x.get_absolute_cost();
        j["maxEvaluations"] = x.get_max_evaluations();
        j["maxIterations"] = x.get_max_iterations();
        j["relativeCost"] = x.get_relative_cost();
    }

    inline void from_json(const json & j, Config& x) {
        x.set_cams(j.at("cams").get<std::map<std::string, V2KSchema>>());
        x.set_solver(j.at("solver").get<Solver>());
        x.set_window(j.at("window").get<double>());
    }

    inline void to_json(json & j, const Config & x) {
        j = json::object();
        j["cams"] = x.get_cams();
        j["solver"] = x.get_solver();
        j["window"] = x.get_window();
    }
}
}
