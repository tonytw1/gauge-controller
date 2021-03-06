package uk.co.eelpieconsulting.monitoring.model;

import org.joda.time.DateTime;

import java.util.List;

public class Metric implements Comparable {

	private final String name, lastValue;
	private final MetricType type;
	private final DateTime date;
	private List<DateTime> changes;

	public Metric(String name, MetricType type, String lastValue, DateTime date, List<DateTime> changes) {
		super();
		this.name = name;
		this.type = type;
		this.lastValue = lastValue;
		this.date = date;
		this.changes = changes;
	}

	public String getName() {
		return name;
	}

	public MetricType getType() {
		return type;
	}

	public String getLastValue() {
		return lastValue;
	}
	
	public DateTime getDate() {
		return date;
	}

	public List<DateTime> getChanges() {
		return changes;
	}

	@Override
	public String toString() {
		return "Metric{" +
				"name='" + name + '\'' +
				", lastValue='" + lastValue + '\'' +
				", type=" + type +
				", date=" + date +
				", changes=" + changes +
				'}';
	}

	@Override
	public int hashCode() {
		final int prime = 31;
		int result = 1;
		result = prime * result + ((name == null) ? 0 : name.hashCode());
		return result;
	}

	@Override
	public boolean equals(Object obj) {
		if (this == obj)
			return true;
		if (obj == null)
			return false;
		if (getClass() != obj.getClass())
			return false;
		Metric other = (Metric) obj;
		if (name == null) {
			if (other.name != null)
				return false;
		} else if (!name.equals(other.name))
			return false;
		return true;
	}

	@Override
	public int compareTo(Object arg0) {
		return name.toLowerCase().compareTo(((Metric) arg0).getName().toLowerCase());
	}
	
}
