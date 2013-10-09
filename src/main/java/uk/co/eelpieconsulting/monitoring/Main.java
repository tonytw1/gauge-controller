package uk.co.eelpieconsulting.monitoring;

import org.apache.velocity.app.VelocityEngine;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.ViewResolver;
import org.springframework.web.servlet.view.velocity.VelocityConfig;
import org.springframework.web.servlet.view.velocity.VelocityConfigurer;
import org.springframework.web.servlet.view.velocity.VelocityViewResolver;

@EnableAutoConfiguration
@ComponentScan
@Configuration
public class Main {

    private static ApplicationContext ctx;
    
	public static void main(String[] args) throws Exception {
    	ctx = SpringApplication.run(Main.class, args);
    }
    
    @Bean
    public VelocityConfig getVelocityConfig() {
    	return new VelocityConfigurer();
    }
    
    @Bean
    public VelocityEngine getVelocityEngine() {
    	return new VelocityEngine();
    }
    
    @Bean
    public ViewResolver getViewResolver() {
    	VelocityViewResolver velocityViewResolver = new VelocityViewResolver();
    	velocityViewResolver.setSuffix(".vm");
    	return velocityViewResolver;
    }
    
}
