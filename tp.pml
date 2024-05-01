#define N 4
#define signal(S) atomic{ S = S+1}

short S=1 //Valor inicial del semáforo en este caso esperamos que llegue a 4 para poder correr la formula de regresión lineal

int results[10];
int m,b;
int ls;

typedef punto {
    int kWh;
    int costo;
};

punto puntos[4];

chan sumX = [1] of {int};
chan sumY = [1] of {int};
chan sumXY = [1] of {int};
chan sumXX = [1] of {int};


active proctype calcularSumX() {
    int subSumX = 0;
    byte i;
    byte n = N;
    for (i : 0 .. N-1) {
        subSumX = puntos[i].kWh +subSumX; 
    }
    sumX[0] = subSumX;
    signal(S);

}

active proctype calcularSumY() {
    int subSumY = 0;
    byte i;
    byte n = N;
    for (i : 0 .. N-1) {
        subSumY = puntos[i].costo +subSumY;
    }
    sumY[0] = subSumY;
    signal(S);

}

active proctype calcularSumXY() {
    int subSumXY = 0;
    byte i;
    byte n = N;
    for (i : 0 .. N-1) {
        subSumXY = subSumXY + (puntos[i].costo * puntos[i].kWh)
    }
    sumXY[0] = subSumXY;
    signal(S);


}

active proctype calcularSumXX() {
    int subSumXX = 0;
    byte i;
    byte n = N;
    for (i : 0 .. N-1) {
        subSumXX = subSumXX + (puntos[i].kWh * puntos[i].kWh)
    }
    sumXX[0] = subSumXX;
    signal(S);



}

active proctype valoresReglineal(){
  
   
    //working
   do ::
        if
        :: (S > 4) -> break;
        fi
    od
   //working


    int subSumX = sumX[0];
    int subSumY = sumY[0];

    int auxSumXY = sumXY[0];
    int auxSumXX = sumXX[0]
    int promX = subSumX / N;
    int promY = subSumY / N;
 


    
    m = ((auxSumXY - N*promX*promY) / ( auxSumXX- N*promX*promY));
    


    b = promY - m*promX;
    b = b 
    signal(S)

}

proctype predictValor(int i){
    int r = m*((300+i)) + b; 
    results[i] = r 
}



init {
    int cont = 0;

    puntos[0].kWh = 150;
    puntos[1].kWh = 200;
    puntos[2].kWh = 250;
    puntos[3].kWh = 300;

    puntos[0].costo = 30;
    puntos[1].costo = 40;
    puntos[2].costo = 50;
    puntos[3].costo = 60;

    
    do ::
        if
        :: (S > 5) -> break;
        fi
    od
    
    
    
    do ::

        

        if
        :: (cont > 5) -> break;
        :: else 
        fi


        

        run predictValor(cont);

        cont = cont + 1;
        
    od

    


}