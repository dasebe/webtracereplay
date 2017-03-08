TARGET = client
OBJS += client.o
LIBS += -lcurl -lpthread  

CXX = g++
CXXFLAGS += -std=c++11
CXXFLAGS += -MMD -MP # dependency tracking flags
CXXFLAGS += -I./
CXXFLAGS += -Wall -Werror 
LDFLAGS += $(LIBS)
all: CXXFLAGS += -O2 # release flags
all:		$(TARGET)

debug: CXXFLAGS += -ggdb  -D_GLIBCXX_DEBUG # debug flags
debug: $(TARGET)

$(TARGET):	$(OBJS)
	$(CXX) $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

%.o: %.c
	$(CXX) $(CXXFLAGS) -c -o $@ $<

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c -o $@ $<

DEPS = $(OBJS:%.o=%.d)
-include $(DEPS)

clean:
	-rm $(TARGET) $(OBJS) $(DEPS)

test:	all	
	./${TARGET} traces/example.json
