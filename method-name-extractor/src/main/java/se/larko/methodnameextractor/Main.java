package se.larko.methodnameextractor;

import antlr4.java8.Java8Lexer;
import antlr4.java8.Java8Parser;
import org.antlr.v4.runtime.ANTLRInputStream;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.tree.ParseTree;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Optional;

public class Main {

    /**
     * Generates the Java8Lexer and Java8Parser classes
     */
    private static void gen() {
        String[] arg0 = { "-visitor", "/home/erik/Code/Privat/bug-stats/method-name-extractor/src/main/java/antlr4/java8/Java8.g4", "-package", "antlr4.java8" };
        org.antlr.v4.Tool.main(arg0);
    }

    /**
     * Usage: Main FILE LINE1 LINE2 LINE3
     * Example:
     *      /home/erik/Code/Privat/bug-stats/method-name-extractor/src/main/java/se/larko/methodnameextractor/Main.java 27 1 48 47
     */
    public static void main(String[] args) throws IOException {
        String fileToParse = args[0];
        LineToMethodRegister register = buildLineToMethodRegister(fileToParse);

        writeResultsToStdOut(args, register);
    }

    private static void writeResultsToStdOut(String[] args, LineToMethodRegister register) {
        System.out.println("[");
        for (int i = 1; i < args.length; i++) {
            System.out.println(String.join(",\n", printLine(register, args[i])));
        }
        System.out.println("]");
    }

    private static String[] printLine(LineToMethodRegister register, String arg) {
        if (arg.contains(" ")) {
            String[] rawLines = arg.split(" ");

            String[] toReturn = new String[rawLines.length];
            for (int i = 0; i < rawLines.length; i++) {
                int line = Integer.parseInt(rawLines[i]);
                toReturn[i] = getMethodOnLineAsJson(register, line);
            }
            return toReturn;
        } else {
            int line = Integer.parseInt(arg);
            return new String[] { getMethodOnLineAsJson(register, line) };
        }
    }

    private static String getMethodOnLineAsJson(LineToMethodRegister register, int line) {
        Optional<LineToMethodRegister.Method> methodOnLine = register.getMethodOnLine(line);
        return LineAndMethodJsonOutput.toJson(line, methodOnLine);
    }

    private static LineToMethodRegister buildLineToMethodRegister(String fileToParse) throws IOException {
        ParseTree tree = getParseTree(fileToParse);

        LineToMethodRegister register = new LineToMethodRegister();
        new Java8ParseTreeVisitor(register).visit(tree);

        return register;
    }

    private static ParseTree getParseTree(String fileToParse) throws IOException {
        ANTLRInputStream input = new ANTLRInputStream(new FileInputStream(fileToParse));

        Lexer lexer = new Java8Lexer(input);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        Java8Parser parser = new Java8Parser(tokens);
        ParseTree tree = parser.compilationUnit();
        return tree;
    }
}
