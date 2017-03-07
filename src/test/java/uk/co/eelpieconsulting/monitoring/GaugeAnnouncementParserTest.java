package uk.co.eelpieconsulting.monitoring;

import org.junit.Test;
import uk.co.eelpieconsulting.monitoring.model.Gauge;
import uk.co.eelpieconsulting.monitoring.model.GaugeType;

import static junit.framework.TestCase.assertNotNull;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;

public class GaugeAnnouncementParserTest {

    private GaugeAnnouncementParser gaugeAnnouncementParser = new GaugeAnnouncementParser();

    @Test
    public void shouldIgnoreMessagesWhichAreNotInGaugeAnnounementFormat() {
        assertNull(gaugeAnnouncementParser.parse("nothing to see here"));
    }

    @Test
    public void shouldParseGaugeFromAnnoucementMessage() {
        Gauge gauge = gaugeAnnouncementParser.parse("lamp:red[1]");
        assertNotNull(gauge);
        assertEquals("red", gauge.getName());
        assertEquals(GaugeType.LAMP, gauge.getType());
        assertEquals(1, gauge.getFsd());
    }

    @Test
    public void shouldIgnoreUnknownGaugeTypes() {
        Gauge gauge = gaugeAnnouncementParser.parse("widget:red[1]");
        assertNull(gauge);
    }

}
