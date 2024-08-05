# kvhammer

todo:
get payloads working.
get writes working 
    Add custom payload.
Get metadata working
figure out how to randomly launch all

handle signals or context.
case/switch.
Make runner function. Think of the data needed to pass i/o.
Channels will work better for the communicating vs sharing of data.
slice of function calls.
verify data existence before patch/update.
directed graph - not yet, but learn.
drop secret mount paths to a map pointer, pick off the reads.
wait for existence checks from vault requests. - pre metadata and patching.
Only dump in to the map for reads if successful return from Vault. - add error handling.
syncmap? Drop a slice and leave it alone for the reads.
set a done or wait group channel. Done channel - give routine and ID, channel of structs, pass in the channel, sends empty struct.
BE SURE YOU CAN STOP IT. - go by example on signals.
nc -ldp udp or tcp, redirect output to dev/nulll.
netcat install in dockerfile though. - needs UDP functionality.