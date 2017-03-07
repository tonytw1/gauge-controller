package uk.co.eelpieconsulting.monitoring;

import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Gauge;
import uk.co.eelpieconsulting.monitoring.model.GaugeType;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Component
public class GaugeAnnouncementParser {

    Pattern pattern = Pattern.compile("(.*?):(.*?)\\[(.*?)\\]");

    public Gauge parse(String message) {
        Matcher matcher = pattern.matcher(message);
        if (matcher.matches()) {
            try {
                GaugeType type = GaugeType.valueOf(matcher.group(1).toUpperCase());
                return new Gauge(matcher.group(2), type, Integer.parseInt(matcher.group(3)));
            } catch (IllegalArgumentException e) {
                return null;
            }
        }
        return null;
    }

}
