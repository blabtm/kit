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
namespace ccd {
    using nlohmann::json;

    #ifndef NLOHMANN_UNTYPED_v2k_ccd_HELPER
    #define NLOHMANN_UNTYPED_v2k_ccd_HELPER
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
        double beta;
        double dispersion;

        public:
        const double & get_beta() const { return beta; }
        double & get_mutable_beta() { return beta; }
        void set_beta(const double & value) { this->beta = value; }

        const double & get_dispersion() const { return dispersion; }
        double & get_mutable_dispersion() { return dispersion; }
        void set_dispersion(const double & value) { this->dispersion = value; }
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

    class Config {
        public:
        Config() = default;
        virtual ~Config() = default;

        private:
        std::map<std::string, V2KSchema> cams;

        public:
        const std::map<std::string, V2KSchema> & get_cams() const { return cams; }
        std::map<std::string, V2KSchema> & get_mutable_cams() { return cams; }
        void set_cams(const std::map<std::string, V2KSchema> & value) { this->cams = value; }
    };
}
}

namespace v2k {
namespace ccd {
    void from_json(const json & j, X & x);
    void to_json(json & j, const X & x);

    void from_json(const json & j, Axes & x);
    void to_json(json & j, const Axes & x);

    void from_json(const json & j, V2KSchema & x);
    void to_json(json & j, const V2KSchema & x);

    void from_json(const json & j, Config & x);
    void to_json(json & j, const Config & x);

    inline void from_json(const json & j, X& x) {
        x.set_beta(j.at("beta").get<double>());
        x.set_dispersion(j.at("dispersion").get<double>());
    }

    inline void to_json(json & j, const X & x) {
        j = json::object();
        j["beta"] = x.get_beta();
        j["dispersion"] = x.get_dispersion();
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

    inline void from_json(const json & j, Config& x) {
        x.set_cams(j.at("cams").get<std::map<std::string, V2KSchema>>());
    }

    inline void to_json(json & j, const Config & x) {
        j = json::object();
        j["cams"] = x.get_cams();
    }
}
}
