# ip-counter

There are two main ideas behind this code.

1) I aim to use less memory than a hashmap and to guarantee O(1) complexity for each lookup / write operation. I achieve this with the help of a table, where each cell corresponds to a certain IP. The table is four-dimensional; it contains ```[256][256][256][32]``` bytes. To save memory, I use every bit of these bytes as a 1-bit variable. I check and set the values of the bits using bit arithmetic.
  
   The function ```threadLoop``` parses an IP string into four ints and uses them to calculate indices of this table. Then it increments the variable that stores the number of non-unique IPs encountered by this thread.

4) I leverage concurrency in order to do as much work as possible in parallel. We can divide all the IP addresses in the file into non-ovelapping categories, for example, all IPs that start with a number less than 32, a number equal or greater to 32 and less than 64, etc. The number of unique IPs will then be the sum of the numbers of unique IPs in all of these categories. We can count this number independently for each category. Moreover, because these categories overlap, all the threads can safely write into the shared table.
  
   I start ```MAX_THREADS``` number of threads, each corresponding to one such category, and create a channel for each of the threads. Then I read the input file using the main goroutine, check which category the current IP belongs to, and send it into the channel belonging to one of the threads. After all threads finish processing the lines from their channels, the overall answer is calculated.

This solution found 1000000000 unique IPs in the 120Gb file; processing the file took 37 minutes.
