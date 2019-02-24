#include "pistache/endpoint.h"
#include <pistache/http_headers.h>
#include <pistache/http.h>

using namespace Pistache;


// block of data
#define BLOCKSIZE 131072  //32768
char dataBig[BLOCKSIZE];

class TestHandler : public Http::Handler {
public:

    HTTP_PROTOTYPE(TestHandler)

    void onRequest(const Http::Request& request, Http::ResponseWriter response) override{
        UNUSED(request);
        // response.headers().add<Header::Server>("lys")
        //      .add<Header::ContentType>(MIME(Text, Plain));
        // auto stream = response.stream(Http::Code::Ok);
        // for (size_t i = 0; i < 10; ++i) {
        //     stream.write(dataBig, BLOCKSIZE);
        // }
        // stream.flush();
        // stream.ends();

        // std::string out = "";
        // for(int n=0; n<6000; n++) {
        //     out += "xxxxxxxxxx";
        // }
        response.send(Http::Code::Ok, dataBig);
    }
};

int main() {

    for (int i=0; i<BLOCKSIZE; i++)
        dataBig[i] = (char)i;

    Pistache::Address addr(Pistache::Ipv4::any(), Pistache::Port(9080));
    auto opts = Pistache::Http::Endpoint::options()
        .threads(1);

    Http::Endpoint server(addr);
    server.init(opts);
    server.setHandler(Http::make_handler<TestHandler>());
    server.serve();

    server.shutdown();
}
