<?xml version="1.0" encoding="utf-8" ?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" >
<xsl:output omit-xml-declaration="yes" indent="yes"/>

<xsl:template match="/report">

<html><head><title></title>
<style type="text/css">
    <!--/*--><![CDATA[/*><!--*/
    body { font-family: monospace ; }
    tr.header { background-color: #C0B0D8; }
    tr.odd { background-color: #D0D0D0; }
    tr.even { background-color: #F0F0F0 ; }
    td.passed { color: green; }
    td.failed { color: red; }
    /*]]>*/-->
</style>
</head>
<body>
    <xsl:apply-templates select="config/property_list"/>
    <xsl:choose>
    <xsl:when test="@type = 'performance'">
        <xsl:apply-templates select="test_list" mode='benchmark'/>
    </xsl:when>
    <xsl:otherwise>
        <xsl:apply-templates select="test_list" mode='auto'/>
    </xsl:otherwise>
    </xsl:choose>
</body>
</html>
</xsl:template>

<!-- CONFIG -->
<xsl:template match="config/property_list">
    <h3>Config</h3>
    <table>
    <tr class="header"><td> Property </td><td> Value </td></tr>
        <xsl:apply-templates select='item' mode='properties'/> 
    </table>
</xsl:template>

<xsl:template match="item" mode='properties'>
    <tr><td><xsl:value-of select="normalize-space(./key)"/></td>
    <td><xsl:value-of select="normalize-space(./value)"/></td></tr>
</xsl:template>


<!-- RESULTS -->
<xsl:template match="test_list" mode='auto'>
    <h3>Results</h3>
    <table>
    <xsl:apply-templates select='item' mode='test_item' />
    </table>
</xsl:template>

<xsl:template match="test_list" mode='benchmark'>
    <h3>Results</h3>
    <table>
    <tr class="header"><td> Name</td><td> Value </td><td> Unit </td><td> Trend </td></tr>
    <xsl:apply-templates select='item' mode='benchmark_item' />
    </table>
</xsl:template>

<xsl:template match='item' mode='benchmark_item'>
    <xsl:variable name="oddeven">
        <xsl:choose>
        <xsl:when test="position() mod 2 = 0">even</xsl:when>
        <xsl:otherwise>odd</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <tr class="{$oddeven}"><td><xsl:value-of select="normalize-space(./name)"/></td>
    <td><xsl:value-of select="normalize-space(./value)"/></td>
    <td><xsl:value-of select="normalize-space(./unit)"/></td>
    <td><xsl:value-of select="normalize-space(./trend)"/></td></tr>
</xsl:template>

<xsl:template match='item' mode='test_item'>
    <xsl:variable name='verdict'><xsl:value-of select='status' /></xsl:variable>
    <xsl:variable name="evenodd">
        <xsl:choose>
        <xsl:when test="position() mod 2 = 0">even</xsl:when>
        <xsl:otherwise>odd</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <tr class="{$evenodd}"><td class='{$verdict}'><xsl:value-of select="normalize-space(./name)"/></td>
    <td class='{$verdict}'><xsl:value-of select="normalize-space(./status)"/></td>
    <td class='{$verdict}'><xsl:value-of select="normalize-space(./message)"/></td></tr>
</xsl:template>

</xsl:stylesheet>

