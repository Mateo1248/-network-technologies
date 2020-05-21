import java.util.Collections;
import java.sql.Time;
import java.util.ArrayList;

class Packet implements Comparable<Packet>{
    Integer id;
    String msg;
    long time;

    Packet(Integer id, String msg) {
        this.id = id;
        this.msg = msg;
    }

    Packet(Integer id, String msg, long time) {
        this.id = id;
        this.msg = msg;
        this.time = time;
    }

    @Override
    public int compareTo(Packet p) {
        return (this.id).compareTo(p.id);
    }

    public static void addToBuff(ArrayList<Packet> buff, Packet p) {
        boolean contains = false;
        for(Packet x : buff) {
            if(x.id == p.id) {
                contains = true;
                break;
            }
        }

        if(!contains) {
            buff.add(p);
            Collections.sort(buff);
        }
    }

    public static int printFromMin(ArrayList<Packet> arrL, int min, String tag) {
        int size = arrL.size();

        while(size != 0 ) {
            Packet p = arrL.get(0);
            
            if(p.id == min) {
                System.out.println(tag + ": " + Integer.toString(p.id) + ": " +  p.msg);
                arrL.remove(0);
                min++;
                size --;
            } else if(p.id < min) {
                arrL.remove(0);
            }
            else {
                break;
            }
            size = arrL.size();
        }
        return min;
    }
}