package se.larko.methodnameextractor;

import antlr4.java8.Java8BaseVisitor;
import antlr4.java8.Java8Parser;

import java.util.Collection;
import java.util.LinkedList;
import java.util.List;

public class Java8ParseTreeVisitor extends Java8BaseVisitor<Void> {

    private final LineToMethodRegister register;

    public Java8ParseTreeVisitor(LineToMethodRegister register) {
        this.register = register;
    }

    @Override
    public Void visitMethodDeclaration(Java8Parser.MethodDeclarationContext ctx) {
        String methodName = extractMethodName(ctx);
        Collection<String> parameterTypes = extractParameterTypes(ctx);

        register.registerMethod(methodName, parameterTypes, ctx.getStart().getLine(), ctx.getStop().getLine());
        //System.out.println("Found method " + methodName + "(" + String.join(" ", parameterTypes) + ") on lines " + ctx.getStart().getLine() + " to " + ctx.getStop().getLine());
        return super.visitMethodDeclaration(ctx);
    }

    private String extractMethodName(Java8Parser.MethodDeclarationContext ctx) {
        return ctx.methodHeader().methodDeclarator().Identifier().getText();
    }

    private Collection<String> extractParameterTypes(Java8Parser.MethodDeclarationContext ctx) {
        List<String> parameterTypes = new LinkedList<>();

        Java8Parser.FormalParameterListContext parameterList = ctx.methodHeader().methodDeclarator().formalParameterList();
        if (parameterList != null) {

            if (parameterList.formalParameters() != null) {
                List<Java8Parser.FormalParameterContext> formalParameterContexts = parameterList.formalParameters().formalParameter();
                for (Java8Parser.FormalParameterContext formalParameterContext : formalParameterContexts) {
                    String parameterType = formalParameterContext.unannType().getText();

                    parameterTypes.add(parameterType);
                }
            }

            String lastParameterType = parameterList.lastFormalParameter().formalParameter().unannType().getText(); // gives the types in public static void main(String[])
            parameterTypes.add(lastParameterType);
        }

        return parameterTypes;
    }
}
