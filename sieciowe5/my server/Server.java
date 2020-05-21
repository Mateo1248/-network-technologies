import com.sun.net.httpserver.Headers;
import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpHandler;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.List;

//Aby połączyć się z serwerem wchodzimy http://localhost:1235/myserver?home 
public class Server {
    public static void main(String[] args) throws IOException {
        //server lokalny na porcie 1235
        HttpServer server = HttpServer.create(new InetSocketAddress(1235), 0);

        //poczatkowy adres to localhost:1235/myserver
        HttpContext context = server.createContext("/myserver"); 

        //handler dla zadan
        context.setHandler(new HttpHandler() {
            @Override
            public void handle(HttpExchange exchange) throws IOException {
                //pobieranie stron
                String response = createHeaderRequest(exchange);
                //wyslanie odp klientowi rcode 200 = operation successful
                exchange.sendResponseHeaders(200, response.getBytes().length);
                OutputStream os = exchange.getResponseBody();
                os.write(response.getBytes());
                exchange.close();
            }
        });
        server.start();
    }

    private static String createHeaderRequest(HttpExchange exchange) {

        StringBuilder response = new StringBuilder(); 
        response.append(handleQuerry(exchange.getRequestURI().getQuery())); // pobieranir metody jaka została wywołana

        response.append("<pre>\n<b>Your header request: \n"); 
        response.append("\nHEADERS:</b>\n");
        Headers requestHeaders = exchange.getRequestHeaders();

        for (String name : requestHeaders.keySet()) {
            List<String> values = requestHeaders.get(name);
            response.append("\n<b>NAME:</b> " + name + " <b>VALUE:</b> ");
            for (String value : values) {
                response.append(value + " ");
            }
            response.append("\n");
        }

        response.append("\n<b>HTTP METHOD:</b>\n");
        response.append(exchange.getRequestMethod());

        response.append("\n\n<b>QUERY:</b>\n");
        response.append(exchange.getRequestURI()).append("</pre></body></html>");


        System.out.println(response);
        return response.toString();
    }

    //ładowanie pliki w zaleznosci od zapytania
    private static String handleQuerry(String query) {
        String context;
        switch (query) { 
            case "first":
                context = loadFile("first.html");
                break;
            case "second":
                context = loadFile("second.html");
                break;
            default:
                context = loadFile("home.html");
        }
        return context;
    }

    private static String loadFile(String file) {

        String text = "no response\n";
        try(BufferedReader br = new BufferedReader(new FileReader("/home/mateusz/Pulpit/sieciowe/sieciowe5/my server/" + file))) {
            StringBuilder sb = new StringBuilder();
            String line;

            while ((line = br.readLine()) != null) {
                sb.append(line);
                sb.append(System.lineSeparator());
            }

            text = sb.toString();
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        return text;
    }
}