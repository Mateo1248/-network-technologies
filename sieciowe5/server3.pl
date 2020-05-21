  use HTTP::Daemon;
  use HTTP::Status;  
  #use IO::File;

  my $d = HTTP::Daemon->new(
           LocalAddr => 'localhost',        #ustawiamy nazwe naszego hosta
           LocalPort => 4321,               #port na którym będzie odbierał nasz serwer
       )|| die;
  
  print "Please contact me at: <URL:", $d->url, ">\n";


  while (my $c = $d->accept) {              #zmienna $c przechowuje wskaznik do połączenia z klientem
      while (my $r = $c->get_request) {     #$r przechowuje wskaznik do żądania
          if ($r->method eq 'GET') {        #jeśli klient dostanie wyśle żądanie 'GET' w odpowiedzi w oknie przeglądarki
                                            #zostanie mu wyświetlona treść pliku index.html
              $file_s= "./index.html";      #index.html - jakis istniejacy plik
              $c->send_file_response($file_s);

          }
          else {
              $c->send_error(RC_FORBIDDEN)
          }

      }
      $c->close;
      undef($c);
  }
