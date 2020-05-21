import java.net.*;
import java.util.ArrayList;

public class Z2Receiver {

    static final int datagramSize=50;
    InetAddress localHost;
    int destinationPort;
    DatagramSocket socket;

    ReceiverThread receiver;

    public Z2Receiver(int myPort, int destPort) throws Exception {
        localHost=InetAddress.getByName("127.0.0.1");
        destinationPort=destPort;
        socket=new DatagramSocket(myPort);
        receiver=new ReceiverThread();    
    }


    class ReceiverThread extends Thread {

        public void run() {
            //<ADDED
            ArrayList<Packet> msgBuff = new ArrayList<>();
            int maxSend = 0;
            //ADDED>
            try {
                while(true) {
                    byte[] data = new byte[datagramSize];
                    DatagramPacket packet = new DatagramPacket(data, datagramSize);
                    socket.receive(packet);
                    Z2Packet p = new Z2Packet(packet.getData());

                    //<ADDED
                    Packet myP = new Packet(p.getIntAt(0), "" + (char)p.data[4]);

                    if(myP.id >= maxSend) {
                        //dodaj do bufora
                        Packet.addToBuff(msgBuff, myP);

                        //sprawdz ciagłosc i wypisz co  sie da
                        maxSend = Packet.printFromMin(msgBuff, maxSend, "Reciver");   
                    }
                    //ADDED>

                    // WYSLANIE POTWIERDZENIA
                    packet.setPort(destinationPort);
                    socket.send(packet);
                }
            }
            catch(Exception e) {
                e.printStackTrace();
                System.out.println("Z2Receiver.ReceiverThread.run: "+e);
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Z2Receiver receiver = new Z2Receiver( Integer.parseInt(args[0]), Integer.parseInt(args[1]));
        receiver.receiver.start();
    }
}