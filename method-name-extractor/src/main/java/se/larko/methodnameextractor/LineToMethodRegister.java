package se.larko.methodnameextractor;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

public class LineToMethodRegister {

    public static class Method {
        private final String name;
        private final Collection<String> parameterTypes;

        public Method(String name, Collection<String> parameterTypes) {
            this.name = name;
            this.parameterTypes = parameterTypes;
        }

        public String getName() {
            return name;
        }

        public Collection<String> getParameterTypes() {
            return parameterTypes;
        }

        @Override
        public String toString() {
            return name + "(" + String.join(" ", parameterTypes) + ")";
        }
    }

    private final Map<Integer, Method> db = new HashMap<>();

    public void registerMethod(String name, Collection<String> parameterTypes, int firstLine, int lastLine) {
        Method method = new Method(name, parameterTypes);
        for (int i = firstLine; i <= lastLine; i++) {
            db.put(i, method);
        }
    }

    public Optional<Method> getMethodOnLine(int line) {
        return Optional.ofNullable(db.get(line));
    }
}
