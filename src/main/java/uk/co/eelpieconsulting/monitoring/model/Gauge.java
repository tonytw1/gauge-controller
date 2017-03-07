package uk.co.eelpieconsulting.monitoring.model;

public class Gauge  {

	private final String name;
	private final GaugeType type;
	private final int fsd;

	public Gauge(String name, GaugeType type, int fsd) {
		this.name = name;
		this.type = type;
		this.fsd = fsd;
	}
	
	public String getName() {
		return name;
	}
	
	public GaugeType getType() {
		return type;
	}

	public int getFsd() {
		return fsd;
	}

	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;

		Gauge gauge = (Gauge) o;

		if (fsd != gauge.fsd) return false;
		if (name != null ? !name.equals(gauge.name) : gauge.name != null) return false;
		return type == gauge.type;

	}

	@Override
	public int hashCode() {
		int result = name != null ? name.hashCode() : 0;
		result = 31 * result + (type != null ? type.hashCode() : 0);
		result = 31 * result + fsd;
		return result;
	}

	@Override
	public String toString() {
		return "Gauge{" +
				"name='" + name + '\'' +
				", type=" + type +
				", fsd=" + fsd +
				'}';
	}

}
