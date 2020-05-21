import java.net.*;
import java.util.ArrayList;

class Z2Sender {
	static final int datagramSize=50;
	static final int sleepTime=500;
	static final int maxPacket=50;
	InetAddress localHost;
	int destinationPort;
	DatagramSocket socket;
	SenderThread sender;
	ReceiverThread receiver;
	//<ADDED
	static final long relayTime = 3000;
	ArrayList<Packet> sended;
	//AEDDED>
	public Z2Sender(int myPort, int destPort) throws Exception {
		localHost = InetAddress.getByName("127.0.0.1");
		destinationPort = destPort;
		socket = new DatagramSocket(myPort);
		sender = new SenderThread();
		receiver = new ReceiverThread();
		sended = new ArrayList<>();
    }

	class SenderThread extends Thread {
		public void run() {
			int i, x;
			try {
				for(i=0; (x=System.in.read()) >= 0 ; i++) {
					Z2Packet p = new Z2Packet(4+1);
					p.setIntAt(i,0);
					p.data[4]= (byte) x;
					DatagramPacket packet = new DatagramPacket(p.data, p.data.length, localHost, destinationPort);

					//<ADDED
					synchronized (sended) {
						socket.send(packet);
						Packet.addToBuff(sended, new Packet(p.getIntAt(0), "" + (char)p.data[4], System.currentTimeMillis()));
					}
					//ADDED>
					
					sleep(sleepTime);

					//<ADDED
					boolean w = false;
					synchronized (sended) {
						if(sended.size() > 0) {
							long curTime = System.currentTimeMillis();
							for(Packet pac : sended) {
								if(curTime - pac.time > relayTime) {
									pac.time = System.currentTimeMillis();
									Z2Packet p2 = new Z2Packet(5);
									p2.setIntAt(pac.id, 0);
									p2.data[4] = (pac.msg).getBytes()[0];
									DatagramPacket packet2 = new DatagramPacket(p2.data, p2.data.length, localHost, destinationPort);
									socket.send(packet2);
									w = true;
								}
							}
						}
					}
					sleep(sleepTime*(w?1:0));
 					//ADDED>
				}
				//<ADDED
				// send by the time you get all confirmation
				while(sended.size() > 0) {
					int ctr = 0;
					synchronized (sended) {
						long curTime = System.currentTimeMillis();
						for(Packet pac : sended) {
							if(curTime - pac.time > relayTime) {
								pac.time = System.currentTimeMillis();
								Z2Packet p2 = new Z2Packet(5);
								p2.setIntAt(pac.id, 0);
								p2.data[4] = (pac.msg).getBytes()[0];
								DatagramPacket packet2 = new DatagramPacket(p2.data, p2.data.length, localHost, destinationPort);
								socket.send(packet2);
							}
						}
					}					
					
					sleep(sleepTime*(ctr+1));
					ctr = 0;
				}

				System.out.println("Get all confirmations.");
				 //ADDED>
			}
			catch(Exception e){
				System.out.println("Z2Sender.SenderThread.run: "+e);
			}
		}
	}


	class ReceiverThread extends Thread {
		public void run() {
			try {
				while(true) {
					byte[] data = new byte[datagramSize];
					DatagramPacket packet = new DatagramPacket(data, datagramSize);
					socket.receive(packet);
					Z2Packet p = new Z2Packet(packet.getData());

					//ADDED
					int packetID = p.getIntAt(0);

					synchronized (sended) {
						int sendedSize = sended.size();
						for (int i =0 ; i < sendedSize ; i++) {
							if(sended.get(i).id == packetID) {
								sended.remove(i);
								i--;
								sendedSize = sended.size();
							} 
						}
					}
					//ADDED

					System.out.println("S:"+p.getIntAt(0)+": "+(char)p.data[4]);
				}
			}
			catch(Exception e) {
				System.out.println("Z2Sender.ReceiverThread.run: "+e);
			}
		}
	}


	public static void main(String[] args) throws Exception {
		Z2Sender sender = new Z2Sender( Integer.parseInt(args[0]),
		Integer.parseInt(args[1]));
		sender.sender.start();
		sender.receiver.start();
    }
}