package se.larko.methodnameextractor;

import com.google.gson.Gson;

import java.util.LinkedList;
import java.util.Optional;

public class LineAndMethodJsonOutput {

    private static final LineToMethodRegister.Method DEFAULT = new LineToMethodRegister.Method("", new LinkedList<>());
    public static final Gson GSON = new Gson();

    public static String toJson(int line, Optional<LineToMethodRegister.Method> method) {
        boolean hasMethod = method.isPresent();
        LineToMethodRegister.Method m = method.orElse(DEFAULT);

        return GSON.toJson(new LineAndMethodJsonOutput(
                line,
                hasMethod,
                m.getName(),
                m.getParameterTypes().toArray(new String[0]))
        );
    }

    private final int line;
    private final boolean hasMethod;
    private final String methodName;
    private final String[] parameterTypes;

    private LineAndMethodJsonOutput(int line, boolean hasMethod, String methodName, String[] parameterTypes) {
        this.line = line;
        this.hasMethod = hasMethod;
        this.methodName = methodName;
        this.parameterTypes = parameterTypes;
    }
}
