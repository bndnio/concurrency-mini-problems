import threading
import queue

class tangle_node:
    verified = False
    approvalCount = 0
    approvalMutex = queue.Queue(maxsize=1)
    workingMutex = queue.Queue(maxsize=2)
    outbound = []

    def __init__(self):
        self.approvalMutex.put(True)
        return None
    
    def Verify(self, addr, approved, cb):
        self.approvalMutex.get()
        if self.approvalCount >= 2 or addr in self.outbound:
            cb.put(False)
            self.approvalMutex.put(True)
            return
        self.approvalMutex.put(True)

        self.workingMutex.put(True)
        verified = self.approved.get()

        if verified:
            self.approvalMutex.get()
            if self.approvalCount >= 2 or addr in self.outbound:
                cb.put(False)
                self.approvalMutex.put(True)
                self.workingMutex.get()
                return
            else:
                self.approvalCount += 1
            
            if self.approvalCount == 2:
                self.verified = True

            self.outbound.append(addr)
            self.approvalMutex.put(True)

        self.workingMutex.get()
        cb.put(True)

l = []
lWriteMutex = queue.Queue(maxsize=1)

def addWord():
    newNode = tangle_node()
    
    verified, i = 0, 0
    while verified < 2:
        comm = queue.Queue(maxsize=1)
        cb = queue.Queue(maxsize=0)
        l[i%len(l)].Verify(newNode, comm, cb)
        #  ** node evaluation work would go here **
        comm.put(True)
        didVerify = cb.get()
        if didVerify:
            verified += 1
        i+= 1

if __name__ == "__main__":
    l.append(tangle_node())
    l.append(tangle_node())
    lWriteMutex.put(True)

    for i in range(2) :
        addWord()

